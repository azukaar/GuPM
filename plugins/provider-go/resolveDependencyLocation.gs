var name = Dependency.name;
var version = Dependency.version;

if(name.match(/^github/)) {
    var realname = name.match(/^github.com\/([\w\.\-\_]+\/[\w\.\-\_]+)/)
    Dependency.url = 'https://github.com/' + realname[1] + '/archive/'+(version || 'master')+'.zip'
} else if(name.match(/^gopkg.in/)) {
    var payload = httpGet("https://"+name)
    versionMatch = payload.match(/https:\/\/github.com\/([\w\.\-\_]+\/[\w\.\-\_]+)\/tree\/([\w\.\-\_]+)/);

    if(!versionMatch) {
        console.error("Couldn't resolve " + name)
        exit()
    }

    packageName = versionMatch[1]
    packageVersion = versionMatch[2]
    
    // Dependency.name = 'github.com/' + packageName;
    Dependency.version = packageVersion;
    Dependency.url = 'https://github.com/' + packageName + '/archive/'+packageVersion+'.zip'
} else {
    Dependency.url = ""
}

Dependency;
