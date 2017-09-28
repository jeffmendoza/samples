# Getting Started with Packer on Google Compute Engine

## Install Packer

Visit the [Packer download page](https://www.packer.io/downloads.html)
and download the packer binary for your OS and CPU architecture. Unzip
and place the binary in a directory that is in your PATH.

Verify that Packer is installed and in your path by running it:

```console
$ packer
Usage: packer [--version] [--help] <command> [<args>]

Available commands are:
    build       build image(s) from template
    fix         fixes templates from old versions of packer
    inspect     see components of a template
    push        push a template and supporting files to a Packer build service
    validate    check that a template is valid
    version     Prints the Packer version

$
```

## Get your GCP credentials

Use one of the two options.

1. Run packer on a GCE VM. Create the VM with the scopes needed to perform packer operations.

2. Use the cloud sdk to set Application Default credentials with `gcloud auth application-default login`

## Review the template

First there is a builder. This instructs packer on how to create the
VM for buidling the image. There are many optional configuration
settings, such as what to name the resulting image

```
{
    "type": "googlecompute",
    "project_id": "<YOUR_PROJECT_ID>",
    "source_image_family": "debian-9",
    "zone": "us-central1-f",
    "ssh_username": "packer"
}
```

Replace `<YOUR_PROJECT_ID>` with your project id.

The second section is the provisioner. This tells packer what to run
to customize the image before capturing it. There are many provisioner
types available in packer, and they are independant of the builder
type. Here we are simply installing Redis.

```json
{
    "type": "shell",
    "inline": [
        "sleep 30",
        "sudo apt-get update",
        "sudo apt-get install -y redis-server"
    ]
}
```

## Run the packer build

```console
# packer build ./packer-gce.json
...
```

Packer will create a temporary VM, run the shell provisioner to
install Redis, then save the resulting image. Check your images.

```console
$ gcloud compute images list --no-standard-images
NAME                                         PROJECT    FAMILY    DEPRECATED  STATUS
packer-1506630388                            my-project                       READY
$ 
```

You now have a customized image with Packer!