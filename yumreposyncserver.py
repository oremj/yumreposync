import argparse
import hashlib
import logging
import os
import time

from subprocess import check_call

from boto import s3, sqs
from boto.sqs.jsonmessage import JSONMessage

logging.basicConfig(level=logging.INFO)

s3_conn = s3.connect_to_region('us-east-1')
sqs_conn = sqs.connect_to_region('us-east-1')


def md5file(filename):
    m = hashlib.md5()
    m.update(open(filename).read())
    return m.hexdigest()


class RepoSync(object):
    def __init__(self, bucket, queue, repodir):
        self.bucket = s3_conn.get_bucket(bucket)

        self.queue = sqs_conn.lookup(queue)
        self.queue.set_message_class(JSONMessage)

        self.repodir = repodir

    def fetch_keys(self, keys):
        for key_name in keys:
            key = self.bucket.get_key(key_name)
            if not key:
                logging.info("Key: %s does not exist." % key_name)
                continue

            outfile = os.path.join(self.repodir, key.name[1:])
            logging.info("Fetching: %s to %s" % (key.name, outfile))
            key.get_contents_to_filename(outfile)

    def loop(self):
        while True:
            m = self.queue.read()
            if not m:
                time.sleep(10)
                continue
            keys = m['keys']
            self.fetch_keys(keys)
            self.update_metadata()
            self.sync_metadata()
            m.delete()

    def sync_metadata(self):
        existing = list(self.bucket.list("repodata"))
        repodata_dir = os.path.join(self.repodir, "repodata")

        synced = set()
        for f in os.listdir(repodata_dir):
            filename = os.path.join(repodata_dir, f)
            keyname = os.path.join("/repodata", f)
            synced.add(keyname[1:])

            key = self.bucket.get_key(keyname)
            if key and md5file(filename) == key.etag.strip('"'):
                continue

            logging.info("Adding: %s" % filename)
            key = self.bucket.new_key(keyname)
            key.set_contents_from_filename(filename)

        for k in existing:
            if k.name not in synced:
                logging.info("Removing: %s" % k.name)
                k.delete()

    def update_metadata(self):
        check_call(['createrepo', '.'], cwd=self.repodir)


parser = argparse.ArgumentParser(description="yumreposync server")
parser.add_argument('bucket', help='S3 bucket name')
parser.add_argument('queue', help='SQS queue name')
parser.add_argument('repodir', help='Local repo directory')

args = parser.parse_args()

RepoSync(args.bucket, args.queue, args.repodir).loop()
