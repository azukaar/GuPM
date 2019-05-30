![dog](./banner.png)

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
-- LINUX or Mac OS
curl -fsSL  https://azukaar.github.io/GuPM/install.sh | sudo bash 
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

## Env Manager

```
g env A=B env
....
A=B
```

## .gupm_rc.gs

equivalent of .bash_rc but written in .gs.
Put at the root of your folder, will be executed every time you execute `g ...` in your folder

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

### New projects

In order to simply bootstrap a new project you can run `g bootstrap` you can also use `b` and add a provider `g b -p npm`

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

### VS Code 

Add this to your `settings.json` to treat .gs file as javascript (temporary fix to plain text)

```
"files.associations": {
    "*.gs": "javascript"
}
```

### Thanks you!
Package Icon made by [smashicons](https://www.smashicons.com/)
Dog Icon made by [Freepik](https://www.freepik.com/)