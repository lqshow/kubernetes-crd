#!/usr/bin/env bash
os=darwin
arch=amd64

# download kubebuilder and extract it to tmp
curl -L https://go.kubebuilder.io/dl/2.3.1/$\{os\}/$\{arch\} | tar -xz -C /tmp/

# move to a long-term location and put it on your path
# (you'll need to set the KUBEBUILDER_ASSETS env var if you put it somewhere else)
sudo mv /tmp/kubebuilder_2.3.1__ /usr/local/kubebuilder
export PATH=/usr/local/Cellar/go/1.14.3/libexec/bin:/Users/linqiong/miniconda3/bin:/Users/linqiong/miniconda3/condabin:/Library/Java/JavaVirtualMachines/jdk-14.0.1.jdk/Contents/Home/bin:/Applications/Sublime Text.app/Contents/SharedSupport/bin:/Users/linqiong/.nvm/versions/node/v8.4.0/bin:/Users/linqiong/.krew/bin:/Users/linqiong/bin:/usr/local/bin:/usr/bin:/bin:/usr/sbin:/sbin:/Applications/Wireshark.app/Contents/MacOS:/Users/linqiong/workspace/app/golang/lib/bin:/usr/local/bin:/usr/local/opt/python@3.8/bin:/opt/scala/current/bin:/opt/hadoop/current/sbin:/opt/hadoop/current/bin:/opt/kafka/current/bin:/opt/spark/current/bin:/Users/linqiong/workspace/app/golang/lib/bin:/usr/local/kubebuilder/bin
