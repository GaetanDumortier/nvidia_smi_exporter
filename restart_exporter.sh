compute_nodes=("compute-01" "compute-02" "compute-03" "compute-04" "compute-05" "compute-06" "compute-07" "compute-08" "compute-09" "compute-10" "compute-11" "compute-12" "compute-13")

for n in "${compute_nodes[@]}";
do
	# Usefull for when deploying new version of the exporter, as the FS is being synced across every node anyway.
        ssh $n 'systemctl restart nvidia_exporter'
done
