# Godspeed
Godspeed - deployment pipeline tool (spike)

## Usage

Example Shell Deployment:

```
go build -o gs .
./gs deploy --config ./examples/shell.yml
```

You should get an output something like this (if there was a failure during deploy):

```
2015/10/22 14:00:32.946356 [DEBUG]		Loading plugin	shell
2015/10/22 14:00:32.946369 [DEBUG]		Setting up Shell
2015/10/22 14:00:32.946372 [INFO]		Performing Deployment...
2015/10/22 14:00:32.946375 [INFO]		Shell Deploying
2015/10/22 14:00:32.946381 [INFO]		 --> Running command: echo "foo command running!"
2015/10/22 14:00:32.946384 [INFO]		 --> Running command: echo "doing something..."
2015/10/22 14:00:32.946387 [INFO]		 --> Running command: echo "done!"
2015/10/22 14:00:32.946391 [INFO]		Shell Deployment complete with errors, rolling back!
2015/10/22 14:00:32.946394 [INFO]		Shell Rolling back
2015/10/22 14:00:32.946397 [INFO]		 --> Running command: echo "Rollback command running!"
2015/10/22 14:00:32.946400 [INFO]		 --> Running command: echo "./panic --string lucky this is a test"
2015/10/22 14:00:32.946403 [INFO]		 --> Running command: echo "rollback done!"
2015/10/22 14:00:32.946406 [INFO]		Shell Rollback complete!
2015/10/22 14:00:32.946408 [DEBUG]		Shell Teardown
```
