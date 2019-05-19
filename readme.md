![dog](./dog.png)

Global Universal Project Manager -- Package manager, cli tool, scripts for all your projects and your system. Whether you are a developer managing dependencies, or a sysadmin looking for your new toolbelt (bye bash!) you are amoung friends.

---

GuPM is born from the fustration of having to give up my habbits whenever I would switch off Javascript and loose NPM (Whether it would be in Ruby, Go, or even non dev situation). GuPM is claiming to take inspiration from the best things from Brew, NPM, Gem, etc... And compile them in a single tool, usable in any situation.

Example commands :

```
g make -- run make.gs
g install mysql -- locally install MysQL in your project
g mysql -u <user> -p -e "select * from schema.table" -- run the local mysql CLI
mysql -u <user> -p -e "select * from schema.table" -- Does NOT work, MySQL has only been installed in your project, not globally (no version clash)
```

How to install : 

```
-- LINUX 
curl -fsSL  https://azukaar.github.io/GuPM/install.sh | sudo bash 

-- MAC OS
curl -fsSL  https://azukaar.github.io/GuPM/install_mac.sh | sudo bash
```

## Dependency Manager

### Make

This command will setup your project by getting dependencies. Adding a -p or --provider argument allow you to specify what provider to use initially.
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

g install npm://react@1 #will save in gupm.json
g install -p npm react@1 #will save in package.json
```

### remove

## Config Manager
## Script Manager

GuPM also allow you to manage your CLI application.

### Install / use CLI

You can install a CLI application locally to a folder / project, and invoke it in the `g` command

```
g i brew://mysql
g mysql ....
```

You can also install them globally with the -g --global flag

```
g i -g brew://mysql
mysql ....
```

### Write GuPM scripts

You can use GuScript to write bash-like files, used for setting up your project, or use it, or anything really.
Think of GuScript as a replacement for your bash scripts.

```
// name_setup.gs

var name = input('What is your name')
echo('Welcome' + name)
saveName(name)
```

GuScript is based on javascript, and therefore allow advanced object/arrays manipulations, function definitions, etc...
Find more details about the available APIs here:

## CI Manager

## Write plugins
List of function : [./functions.md](Functions)

### Provider

Allow you to install / add from a repo

List of hooks : [./hooklist.md](Hooklist)

Dog Icon made by [Freepik](https://www.freepik.com/)