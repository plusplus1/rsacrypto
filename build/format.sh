export GOPROXY="http://goproxy.apps.intra.yongqianbao.com"

for s in $(git status -s |grep -e '.go$'|grep -v '^D ' |awk '{print $NF}'); do 
    echo "go fmt $s";
    go fmt $s;
done


