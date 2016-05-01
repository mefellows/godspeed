// Check if cluster created...
aws ecs describe-clusters --clusters matts-cluster --region $AWS_REGION

  //  if not, create (need cfn template)
  aws ecs create-cluster --cluster-name "matts-cluster"

// Wait for cluster to become ready
aws ecs describe-clusters --clusters matts-cluster --region $AWS_REGION

aws ecs register-task-definition --family godspeed-hack  --cli-input-json file://./examples/test.json  --region $AWS_REGION

// Grab Task Definition ARN

// Run task

aws ecs create-service  --cluster matts-cluster --service-name pact-broker --task-definition arn:aws:ecs:ap-southeast-2:773592622512:task-definition/godspeed-hack:1 --region ap-southeast-2 --load-balancers loadBalancerName=godspeed-hack,containerName=pact_broker,containerPort=80 --desired-count 1 --role ecsServiceRole


//aws ecs run-task  --cluster matts-cluster --task-definition arn:aws:ecs:ap-southeast-2:773592622512:task-definition/godspeed-hack:1 --region ap-southeast-2

// Grab Task ARN

// Wait for Task done
aws ecs describe-tasks --tasks "arn:aws:ecs:ap-southeast-2:773592622512:task/369a704f-d21d-406c-ab84-b2e767dca799" --cluster matts-cluster --region $AWS_REGION

// Smoke test!
