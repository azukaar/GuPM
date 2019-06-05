![dog](./banner.png)

Global Universal Project Manager -- Package manager, CLI tool, and scripts for all your projects and your system. Whether you are a developer managing dependencies, or a sysadmin looking for your new toolbelt (bye bash!) you are among friends.

 * ⏱**Fast**. Written in native code, with real multi-threading
 * 👓**Smart**. Memory efficient solution using hard-link, which do not duplicate dependencies across project
 * 🌍**Global**. Windows, Mac and Linux compatibility
 * 🌈**Universal**. Usable in any kind of project (Ruby, JS, Go, C, Python, etc...)
 * 👗**Customizable**. Flexible plugin system: make GuPM your own
 * 👝**Future Proof**. Let's make this the last PM you will ever need.

This idea is born from the frustration of having to give up my habits whenever I would switch off Javascript and lose NPM (Whether it would be in Ruby, Go, or even situations outside of coding). GuPM is claiming to take inspiration from the best things in Brew, NPM, Gem, etc... And compile them in a single tool, usable in any project.

 * 📦**Packages Manager**. Install packages from any repository and manage your project's dependency in a seamless way
 * 🖥**CLI Manager**. Install and use CLI tools in a flexible way without conflicts
 * 🚏**Scripting language**. Bye bye Bash, GuPM is bundled with GuScript, allowing you to build cross platform scripts for your project
 * 🐙**Packed with features**. Manage configs, environment variables, CI, and more!
 * 🔥**Even more to come!** See : XXX for the roadmap of feature. You are welcomed to contribute!

---

Here's an example of a workflow using GuPM:

```
g make -- This command will install any dependency your project has
g install mysql -- locally install MySQL in your project
g mysql -u <user> -p -e "select * from schema.table" -- run the local mysql CLI
mysql -u <user> -p -e "select * from schema.table" -- Does NOT work, MySQL has only been installed in your project, not globally (=> no version clash between projects)
```

How to install GuPM : 

```
-- LINUX or Mac OS
curl -fsSL  https://azukaar.github.io/GuPM/install.sh | sudo bash 
```

```
-- Windows
Simply execute: https://azukaar.github.io/GuPM/windows_install.exe
```

## Dependency Manager

### Make

This command will set up your project by getting dependencies. Adding a -p or --provider argument allows you to specify what provider to use initially.
Please note you do NOT need to install npm / gem / whatever to use their corresponding provider, GuPM implement everything itself.

```
g make
g make -p npm
```

### Install

```
-- use default repo

g install mysql
g i mysql
g i node:react

-- use brew

g install brew://mysql
g install -p brew mysql

-- use NPM

g install npm://react@1 -- will save in gupm.json
g install -p npm react@1 -- will save in package.json
```

### remove

```
g remove myPackage
```

### GuPM management

GuPM can be managed using :

```
g self upgrade
g self uninstall
g cache check
g cache clear
```


### Env Manager

```
g env A=B env
....
A=B
```

### .gupm_rc.gs

equivalent of .bash_rc but written in .gs.
Put at the root of your folder, will be executed every time you execute `g ...` in your folder
An example of usage is to setup basic env:

```
env("PATH", pwd())
```

### Install / use CLI

GuPM also allows you to manage your CLI application.
You can install a CLI application locally to a folder/project, and invoke it in the `g` command

```
g i brew://mysql
g mysql ....
```

You can also install them globally with the -g --global flag

```
g i -g brew://mysql
mysql ....
```

### New projects

In order to simply bootstrap a new project you can run `g bootstrap` you can also use `b` and add a provider `g b -p npm`

### Write GuPM scripts

You can use GuScript to write bash-like files, used for setting up your project, use it, or anything literally.
Think of GuScript as a replacement for your bash scripts.

```
// name_setup.gs

var name = input('What is your name')
echo('Welcome' + name)
saveName(name)
```

GuScript is based on javascript, and therefore allow advanced object/arrays manipulations, function definitions, etc...
Find more details about the available APIs here: [./functions.md](Functions)

### CI Manager

## Write plugins

Please find here the documentation to get you started on writing plugins for GuPM

Here's a list of hooks for you to override GuPM's behaviour in your plugin : [./hooklist.md](Hooklist)

List of function available in GS : [./functions.md](Functions)

## VS Code 

Add this to your `settings.json` to treat .gs file as javascript (temporary fix to plain text)

```
"files.associations": {
    "*.gs": "javascript"
}
```

### Thanks!
Package Icon made by [smashicons](https://www.smashicons.com/)
Dog Icon made by [Freepik](https://www.freepik.com/)