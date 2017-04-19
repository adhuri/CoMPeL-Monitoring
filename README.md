# CoMPeL-Monitoring
CoMPeL is a framework which  Monitors resource utilization of a container, predicts resource utilization and does live migration of containers to achieve efficient resource utilization. This Repository consists of monitoring module.


### To run and deploy

#### Build Go Binaries

```chmod +x build.sh```

```./build.sh```

#### Deploy using ansible

```ansible-playbook -i hosts deploy.yml ```



### To release on Github

#### Check latest tag
```git tag```

eg - 1.0.2

#### Use next tag for commit 
```git tag -a 1.0.3 -m " Comment about the next release"```

#### Push using tags for release
 
```git push origin master --tags ```

#### Screenshot for grafana
 
![alt text][grafana]

[grafana]: https://github.com/adhuri/Compel-Monitoring/blob/master/grafana-dashboard/grafana_screenshot.png "Grafana"
