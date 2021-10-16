nohup ./etcd --name s1 \
--listen-client-urls http://0.0.0.0:23791 \
--listen-peer-urls http://0.0.0.0:23801 \
--advertise-client-urls http://0.0.0.0:23791 \
--initial-advertise-peer-urls http://0.0.0.0:23801 \
--initial-cluster s1=http://0.0.0.0:23801,s2=http://0.0.0.0:23802,s3=http://0.0.0.0:23803 \
--initial-cluster-token tkn \
--initial-cluster-state new \
--log-level info \
--logger zap \
--log-outputs stderr \
> s1.out 2>&1 &

nohup ./etcd --name s2 \
--listen-client-urls http://0.0.0.0:23792 \
--listen-peer-urls http://0.0.0.0:23802 \
--advertise-client-urls http://0.0.0.0:23792 \
--initial-advertise-peer-urls http://0.0.0.0:23802 \
--initial-cluster s1=http://0.0.0.0:23801,s2=http://0.0.0.0:23802,s3=http://0.0.0.0:23803 \
--initial-cluster-token tkn \
--initial-cluster-state new \
--log-level info \
--logger zap \
--log-outputs stderr \
> s2.out 2>&1 &

nohup ./etcd --name s3 \
--listen-client-urls http://0.0.0.0:23793 \
--advertise-client-urls http://0.0.0.0:23793 \
--listen-peer-urls http://0.0.0.0:23803 \
--initial-advertise-peer-urls http://0.0.0.0:23803 \
--initial-cluster s1=http://0.0.0.0:23801,s2=http://0.0.0.0:23802,s3=http://0.0.0.0:23803 \
--initial-cluster-token tkn \
--initial-cluster-state new \
--log-level info \
--logger zap \
--log-outputs stderr \
> s3.out 2>&1 &
