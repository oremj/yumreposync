import argparse
import os

from boto import s3, sqs
from boto.sqs.jsonmessage import JSONMessage


sqs_conn = sqs.connect_to_region('us-east-1')
s3_conn = s3.connect_to_region('us-east-1')


def write_files(bucket, files):
    bucket = s3_conn.get_bucket(bucket)

    keys = []
    for f in files:
        key_name = os.path.join("/", os.path.basename(f))
        key = bucket.new_key(key_name)
        key.set_contents_from_filename(f)
        keys.append(key.name)

    return keys


def notify_server(queue, keys):
    q = sqs_conn.lookup(queue)
    m = JSONMessage()
    m['keys'] = keys
    q.write(m)


parser = argparse.ArgumentParser(description="yumreposync client")
parser.add_argument('bucket', help='S3 bucket name')
parser.add_argument('queue', help='SQS queue name')
parser.add_argument('files', nargs='+', help="files to publish")

args = parser.parse_args()

keys = write_files(args.bucket, args.files)
notify_server(args.queue, keys)
