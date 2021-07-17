#!/usr/bin/python

import argparse
import boto3
import csv

sagemaker_featurestore_runtime_client = boto3.Session().client(
  service_name="sagemaker-featurestore-runtime")

# Initialize the parser.
parser = argparse.ArgumentParser()
parser.add_argument("-i", "--input_file", help = "Path to a csv file containing data for ingestion.")
parser.add_argument("-fg", "--feature_group_name", help = "Name of the feature group to write data to.")

# Read arguments from the command line.
args = parser.parse_args()

# Write records from the csv file to s3.
with open(args.input_file) as file_handle:
  for row in csv.DictReader(file_handle, skipinitialspace=True):
    record=[]
    for featureName, valueAsString in row.items():
      record.append({
          'FeatureName':featureName,
          'ValueAsString':valueAsString
      })
    sagemaker_featurestore_runtime_client.put_record(
      FeatureGroupName=args.feature_group_name,
      Record=record)       
