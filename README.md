# dotfiles

My way of handling dotfiles


## Rationale

There is as many ways of handling dotfiles as there is Linux and MacOS X users out there I guess. Make sure you pick the way that suits your needs and style best. [This list](https://dotfiles.github.io/) can be a great inspriation.

In my experience build steps often save lives. The customization options they offer at a defined hook point before production use help me a lot. This is why I took that approach with my dotfiles too.

When you call `./dotfiles.linux64` your dotfiles will be generated for your current system. This way nothing unused clutters your system plus you have the chance to customize the files **on deploy**. Since I am using [Go](https://golang.org) and their awesome [text templating package](https://golang.org/pkg/text/template/), this step gains great flexibility:
* include files in dotfiles that do not inherently support it (e.g. [.ssh/config](http://superuser.com/a/247572))
* change files based on options set in a [simple json file](#json) (yet to be implemented)
* all the other nice features of a real [templating language](http://gohugo.io/templates/go-templates/) (I definitely prefer this over bash trickery)

Yes, I am well aware this is not far from what the great guys at [Ansible](http://www.ansible.com) already built, but Ansible is more generic and thus inherently more complex and requires more setup for me. It still might be the tool of your choice for the task!


## Options

* **src** directory is where you place your custom dotfiles. It is allowed to have subdirectories so you can deploy deep in your home directory. The file named `src/local/.cacheconfig` will be deployed to `~/local/.cacheconfig`. Files with a filename prefix of `._` are not going to be deployed, so those filenames are suitable for files to be included via the template mechanism.
* **backup** directory is used to backup files replaced while deploying the dotfiles you supplied.


## Usage



## Dev Setup
