package strategy

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/mefellows/godspeed/log"
	"github.com/mefellows/plugo/plugo"
	"math/rand"
	"time"
)

type ECSDeploymentStrategy struct {
	ClusterName       string `mapstructure:"cluster_name"`
	Application       string
	Region            string `default:"ap-southeast-2"`
	TaskDefinition    string `mapstructure:"task_definition"`
	ecs               *ecs.ECS
	ELB               string
	taskDefinitionARN string
}

func init() {
	rand.Seed(time.Now().UTC().UnixNano())

	plugo.PluginFactories.Register(func() (interface{}, error) {
		return &ECSDeploymentStrategy{}, nil
	}, "ecs")
}

func (s *ECSDeploymentStrategy) Setup() {
	log.Debug("Setting up ECS ")

	s.ecs = ecs.New(nil)
}

func (s *ECSDeploymentStrategy) createCommand(command string) {
}

func (s *ECSDeploymentStrategy) checkCluster(cluster string) {

	params := &ecs.DescribeClustersInput{
		Clusters: []*string{
			aws.String(cluster),
		},
	}
	resp, err := s.ecs.DescribeClusters(params)

	if err != nil {
		log.Fatalf("%v", err)
		return
	}

	if len(resp.Clusters) == 0 {
		log.Fatalf("Cluster '%s' does not exist", cluster)
	}

	if *resp.Clusters[0].Status != "ACTIVE" {
		log.Info("Cluster not yet active, waiting...")
	}
}

func (s *ECSDeploymentStrategy) registerTask(taskJson string) {
	params := &ecs.RegisterTaskDefinitionInput{
		ContainerDefinitions: []*ecs.ContainerDefinition{ // Required
			{
				Essential: aws.Bool(true),
				Image:     aws.String("amazon/amazon-ecs-sample"),
				Name:      aws.String("godspeedweb"),
				Memory:    aws.Int64(512),
				PortMappings: []*ecs.PortMapping{
					{ // Required
						ContainerPort: aws.Int64(80),
						HostPort:      aws.Int64(80),
						Protocol:      aws.String("tcp"),
					},
				},
			},
		},
		Family: aws.String(s.Application),
	}
	resp, err := s.ecs.RegisterTaskDefinition(params)

	if err != nil {
		log.Fatalf("%v", err)
	}

	s.taskDefinitionARN = *resp.TaskDefinition.TaskDefinitionArn
	log.Info("Task %s created", s.taskDefinitionARN)
}

func (s *ECSDeploymentStrategy) createService() {
	params := &ecs.CreateServiceInput{
		DesiredCount:   aws.Int64(1),       // Required
		ServiceName:    aws.String("demo"), // Required
		TaskDefinition: aws.String(s.taskDefinitionARN),
		Cluster:        aws.String(s.ClusterName),
		LoadBalancers: []*ecs.LoadBalancer{
			{ // Required
				ContainerName:    aws.String("godspeedweb"),
				ContainerPort:    aws.Int64(80),
				LoadBalancerName: aws.String("godspeed-hack"),
			},
		},
		Role: aws.String("ecsServiceRole"),
	}
	_, err := s.ecs.CreateService(params)

	if err != nil {
		log.Fatal(err.Error())
		return
	}
}

func (s *ECSDeploymentStrategy) Deploy() error {
	log.Info("Deploying to ECS")

	s.checkCluster(s.ClusterName)
	s.registerTask(s.TaskDefinition)
	s.createService()

	return nil
}

func (s *ECSDeploymentStrategy) Rollback() error {
	log.Info("ECS Rolling back")
	return nil
}

func (s *ECSDeploymentStrategy) Teardown() {
	log.Debug("ECS Teardown")
}
