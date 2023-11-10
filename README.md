# :warning:
Code in this repository has been extracted from deprecated 
[Cloud integrations](https://github.com/docker/compose-cli). 
It is **not** intended to become a Docker product, but to be (maybe?) adopted by
community and moved to another github organization.

# Compose ECS

This CLI tool makes it easy to run Docker Compose applications on [ECS](https://aws.amazon.com/ecs)


# Usage

Starting with a standard [Compose file](https://compose-spec.io/) - with some restrictions:
```yaml
services:
  jenkins:
    image: jenkins/jenkins:lts
    ports:
      - 8080:8080
```

run `compose-ecs up` and watch CloudFormation creating resources matching the Compose application model:
```
$ compose-ecs up
[+] Running 10/12
 ⠴ demo                          CreateInProgress User Initiated                                                                       145.6s
 ⠿ LogGroup                      CreateComplete                                                                                          2.0s
 ⠿ JenkinsTCP8080TargetGroup     CreateComplete                                                                                         17.0s
 ⠿ CloudMap                      CreateComplete                                                                                         46.1s
 ⠿ DefaultNetwork                CreateComplete                                                                                          6.1s
 ⠴ LoadBalancer                  CreateComplete                                                                                        142.6s
 ⠿ Cluster                       CreateComplete                                                                                          5.1s
 ⠿ JenkinsTaskExecutionRole      CreateComplete                                                                                         19.1s
 ⠿ Default8080Ingress            CreateComplete                                                                                          1.0s
 ⠿ DefaultNetworkIngress         CreateComplete                                                                                          1.0s
 ⠿ JenkinsTaskDefinition         CreateComplete                                                                                          1.0s
 ⠿ JenkinsServiceDiscoveryEntry  CreateComplete                                                                                          1.0s
 ⠋ JenkinsService                CreateInProgress Resource creation Initiated                                                            6.0s

```

In case of a deployment error, check the CloudFormation console for first failure event(s) and reason


Access service logs using `compose-ecs logs`like you would do with docker compose:
```
$compose-ecs logs
jenkins  | Running from: /usr/share/jenkins/jenkins.war
jenkins  | webroot: /var/jenkins_home/war
... 
jenkins  | 2023-04-21 09:47:38.075+0000 [id=29]	INFO	jenkins.install.SetupWizard#init: 
jenkins  | *************************************************************
jenkins  | *************************************************************
jenkins  | *************************************************************
jenkins  | Jenkins initial setup is required. An admin user has been created and a password generated.
jenkins  | Please use the following password to proceed to installation:
jenkins  | 47e2e8c8e9e74e5f85d99b56e794e95d
jenkins  | This may also be found at: /var/jenkins_home/secrets/initialAdminPassword
jenkins  | *************************************************************
jenkins  | *************************************************************
jenkins  | *************************************************************
```

Use `compose-ecs ps` to retrieve public URL for exposed service:
```
$ compose-ecs ps
NAME                                         COMMAND             SERVICE             STATUS              PORTS
task/demo/97db280b80d1407abe2c7e74de8944e5   ""                  jenkins             Running             demo-LoadBa-1V9BXV1VRS6IP-f595d8e2cf1df3d6.elb.eu-west-3.amazonaws.com:8080:8080->8080/tcp
```

Enjoy service running on AWS ... and eventually run `compose-ecs down` to cleanup all resources:
```
$ compose-ecs down
[+] Running 2/4
 ⠋ demo                   DeleteInProgress User Initiated                                                                               45.1s
 ⠋ JenkinsService         DeleteInProgress                                                                                              44.1s
 ⠿ DefaultNetworkIngress  DeleteComplete                                                                                                 0.0s
 ⠿ Default8080Ingress     DeleteComplete                                                                                                 0.0s
...
```

If you want to review or tweak the applied CloudFormation template, run `compose-ecs convert`:
```
$ compose-ecs convert
(...)
  LoadBalancer:
    Properties:
      LoadBalancerAttributes:
      - Key: load_balancing.cross_zone.enabled
        Value: "true"
      Scheme: internet-facing
      Subnets:
      - subnet-xxx
      - subnet-yyy
      - subnet-zzz
      Tags:
      - Key: com.docker.compose.project
        Value: demo
      Type: network
    Type: AWS::ElasticLoadBalancingV2::LoadBalancer
  LogGroup:
    Properties:
      LogGroupName: /docker-compose/demo
    Type: AWS::Logs::LogGroup
```


Please create [issues](https://github.com/docker/compose-ecs/issues) to leave feedback.

