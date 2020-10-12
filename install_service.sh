compute_nodes=("compute-01" "compute-02" "compute-03" "compute-04" "compute-05" "compute-06" "compute-07" "compute-08" "compute-09" "compute-10" "compute-11" "compute-12" "compute-13")

for n in "${compute_nodes[@]}";
do
        ssh $n 'cp /cluster/monitoring/exporters/nvidia_exporter/nvidia_exporter.service /etc/systemd/system/nvidia_exporter.service'
        ssh $n 'systemctl start nvidia_exporter && systemctl enable nvidia_exporter'
done
