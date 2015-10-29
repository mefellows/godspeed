package strategy

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/mefellows/godspeed/log"
	"github.com/mefellows/plugo/plugo"
	"math/rand"
	"time"
)

type ECSDeploymentStrategy struct {
	ecs               *ecs.ECS
	taskDefinitionARN string

	// Config Items
	ClusterName    string `mapstructure:"cluster_name"`
	Application    string
	Region         string `default:"ap-southeast-2"`
	TaskDefinition string `mapstructure:"task_definition"`
	ELB            string
	ElbId          string `mapstructure:"elb_id"`
	Containers     []ContainerDefinition
	Service        Service
}

type Service struct {
	Name        string
	Application string
}

type ContainerDefinition struct {
	Name         string
	Links        []string
	Image        string
	Essential    bool
	PortMappings map[int]int `mapstructure:"port_mappings"`
	Memory       int
	Environment  map[string]string
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

	definitions := make([]*ecs.ContainerDefinition, len(s.Containers))

	for i, d := range s.Containers {
		definitions[i] = &ecs.ContainerDefinition{
			Essential: aws.Bool(d.Essential),
			Image:     aws.String(d.Image),
			Name:      aws.String(d.Name),
			Memory:    aws.Int64(int64(d.Memory)),
			PortMappings: []*ecs.PortMapping{
				{
					ContainerPort: aws.Int64(80),
					HostPort:      aws.Int64(80),
					Protocol:      aws.String("tcp"),
				},
			},
		}
	}

	params := &ecs.RegisterTaskDefinitionInput{
		ContainerDefinitions: definitions,
		Family:               aws.String(s.Application),
	}

	resp, err := s.ecs.RegisterTaskDefinition(params)

	if err != nil {
		log.Fatalf("%v", err)
	}

	s.taskDefinitionARN = *resp.TaskDefinition.TaskDefinitionArn
	log.Info("Task %s created", s.taskDefinitionARN)
}

func (s *ECSDeploymentStrategy) serviceExists(service string) bool {
	params := &ecs.DescribeServicesInput{
		Services: []*string{
			aws.String(service),
		},
		Cluster: aws.String(s.ClusterName),
	}
	resp, err := s.ecs.DescribeServices(params)

	if err != nil {
		log.Fatal(err.Error())
		return false
	}

	if len(resp.Services) > 0 && *resp.Services[0].Status == "ACTIVE" {
		return true

	}
	return false

}

func (s *ECSDeploymentStrategy) createOrUpdateService() {
	exists := s.serviceExists(s.Service.Name)

	// If service does not exist...
	if !exists {
		log.Info(fmt.Sprintf("Service %s not created, creating...", s.Service.Name))
		params := &ecs.CreateServiceInput{
			DesiredCount:   aws.Int64(1),               // Required
			ServiceName:    aws.String(s.Service.Name), // Required
			TaskDefinition: aws.String(s.taskDefinitionARN),
			Cluster:        aws.String(s.ClusterName),
			LoadBalancers: []*ecs.LoadBalancer{
				{ // Required
					ContainerName:    aws.String(s.Service.Application),
					ContainerPort:    aws.Int64(80),
					LoadBalancerName: aws.String(s.ElbId),
				},
			},
			Role: aws.String("ecsServiceRole"),
		}
		_, err := s.ecs.CreateService(params)

		if err != nil {
			log.Fatal(err.Error())
			return
		}

	} else {
		log.Info(fmt.Sprintf("Service %s already created, updating...", s.Service.Name))
		params := &ecs.UpdateServiceInput{
			DesiredCount:   aws.Int64(1),               // Required
			Service:        aws.String(s.Service.Name), // Required
			TaskDefinition: aws.String(s.taskDefinitionARN),
			Cluster:        aws.String(s.ClusterName),
		}
		_, err := s.ecs.UpdateService(params)

		if err != nil {
			log.Fatal(err.Error())
			return
		}

	}
}

func (s *ECSDeploymentStrategy) Deploy() error {
	log.Info("Deploying to ECS")

	s.checkCluster(s.ClusterName)
	s.registerTask(s.TaskDefinition)
	s.createOrUpdateService()

	return nil
}

func (s *ECSDeploymentStrategy) Rollback() error {
	log.Info("ECS Rolling back")
	return nil
}

func (s *ECSDeploymentStrategy) Teardown() {
	log.Debug("ECS Teardown")
}
