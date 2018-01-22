sudo apt-get update
sudo apt-get install apache2 -y
sudo a2ensite default-ssl
sudo a2enmod ssl
sudo service apache2 restart
sudo cat > /var/www/html/index.html << EOF
<!doctype html><html><body>
<h1>Network Status for Shared VPC Terraform/Google Cloud Example</h1>
<h2>VM 1</h2>
<h3>Internal IP: $VM1_INT_IP</h3>
<pre>
$ ping -c 4 -W 1 $VM1_INT_IP
$(ping -c 4 -W 1 $VM1_INT_IP)

$ curl $VM1_INT_IP
$(curl $VM1_INT_IP)
</pre>
<h3>External IP: $VM1_EXT_IP</h3>
<pre>
$ ping -c 4 -W 1 $VM1_EXT_IP
$(ping -c 4 -W 1 $VM1_EXT_IP)

$ curl $VM1_EXT_IP
$(curl $VM1_EXT_IP)
</pre>
<h2>Standalone VM</h2>
<h3>Internal IP: $ST_VM_INT_IP</h3>
$(if [ $ST_VM_INT_IP = $VM1_INT_IP ]; then echo "<h4>Same internal IP as VM1</h4>"; fi)
<pre>
$ ping -c 4 -W 1 $ST_VM_INT_IP
$(ping -c 4 -W 1 $ST_VM_INT_IP)

$ curl $ST_VM_INT_IP
$(curl $ST_VM_INT_IP)
</pre>
<h3>External IP: $ST_VM_EXT_IP</h3>
<pre>
$ ping -c 4 -W 1 $ST_VM_EXT_IP
$(ping -c 4 -W 1 $ST_VM_EXT_IP)

$ curl $ST_VM_EXT_IP
$(curl $ST_VM_EXT_IP)
</pre>
</body></html>

EOF
