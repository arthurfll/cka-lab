kubectl autoscale deployment {{deploy-name}} --min=1 --max=5 --cpu=80 -n {{namespace}}


