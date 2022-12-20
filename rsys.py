#!/usr/bin/env python3

import argparse
import paramiko

def main():
  # Parse command line arguments
  parser = argparse.ArgumentParser()
  parser.add_argument("credentials_file", help="File containing the SSH credentials (username and password) for the remote systems")
  parser.add_argument("hosts_file", help="File containing the hostnames or IP addresses of the remote systems to connect to")
  args = parser.parse_args()

  # Read the credentials from the credentials file
  with open(args.credentials_file, "r") as f:
    username = f.readline().strip()
    password = f.readline().strip()

  # Read the list of hostnames or IP addresses from the hosts file
  with open(args.hosts_file, "r") as f:
    hosts = f.readlines()

  # Iterate over the hostnames or IP addresses
  for host in hosts:
    host = host.strip()

    # Establish an SSH connection to the remote system
    ssh = paramiko.SSHClient()
    ssh.set_missing_host_key_policy(paramiko.AutoAddPolicy())
    ssh.connect(host, username=username, password=password)

    # Search for private SSH keys on the remote system
    stdin, stdout, stderr = ssh.exec_command("find / -name 'id_rsa' 2>/dev/null")
    private_keys = stdout.readlines()

    # Save the output for this remote system to a file
    with open(f"{host}_private_keys.txt", "w") as f:
      for key in private_keys:
        f.write(key)

    # Close the SSH connection
    ssh.close()

if __name__ == "__main__":
  main()
