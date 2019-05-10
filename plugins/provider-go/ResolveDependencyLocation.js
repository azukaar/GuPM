var name = Dependency.name;
var version = Dependency.version;


if(name.match(/^github/)) {
    Dependency.url = 'https://' + name + '/archive/master.zip'
} else if(name.match(/^gopkg.in/)) {
    console.log(98765)
    versionMatch = name.match(/([\w\.\-\_]+\/[\w\.\-\_]+)\.v(\d+\.?[\d+]?\.?[\d+]?)$/);
    packageName = versionMatch[1].replace(/^gopkg\.in/, 'go-yaml')
    packageVersion = versionMatch[2]
    var url = 'https://api.github.com/repos/' + packageName + '/branches';

    var payload = httpGet(url);
    console.log(payload);
}


// https://github.com/src-d/go-git/archive/6e931e4fdefa202c76242109453447182ae16444.zip

Dependency;
