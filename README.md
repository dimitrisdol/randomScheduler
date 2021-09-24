# Random Scheduler

This is a scheduler out-of-tree plugin in order to randomly asssign different types of applications to server sockets.

It implements a typical random algorithm of matching applciations in pairs.

In order to run:

   1. cd to root directory and type *make*
   2. wait for the build to happen
   3. *optional* docker push the docker image to your hub if you want
   4. docker container exec and transfer the random-config.yml file to /etc/kubernetes and the kube-scheduler.yml file to /etc/kubernetes/manifests folder.
   5. validate the plugin with the proposed Pod template.
