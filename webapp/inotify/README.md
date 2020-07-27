docker build -t inotify-img .
docker run -d -p 8090:8090 --name bindmount -v /Users/vivek/temp:/mnt/ --name go-app-container inotify-img 


From Kubernetes
---------------
kubectl create -f deployment.yml

kubectl expose deployment inotify-app --type=NodePort --name=inotify-app-svc --target-port=8090

kubectl get svc
minikube ip
