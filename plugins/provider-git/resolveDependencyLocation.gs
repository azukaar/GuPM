var name = Dependency.name;
var version = Dependency.version;

Dependency.url = 'https://' + name + '/archive/master.zip'

if(Dependency.version === "*.*.*") {
  Dependency.version = "master"
}

// https://github.com/src-d/go-git/archive/6e931e4fdefa202c76242109453447182ae16444.zip

Dependency;
