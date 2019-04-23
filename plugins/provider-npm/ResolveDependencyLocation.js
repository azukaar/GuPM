var name = Dependency.name;
var version = Dependency.version;
var finalVersion;

var payload = httpGet('https://registry.npmjs.org/'+name);
var versionList = Object.keys(payload.versions);

// test tags
if(version.match(/^\d*_*[a-zA-Z]+[\w-_]*[\d\w_]*$/)) {
    finalVersion = payload['dist-tags'][version];
}

// exact version
else if (version.match(/^\d+\.\d+\.\d+/) && !version.match(/\sx/)) {
    finalVersion = version;
}

// test ranges
else {
    finalVersion = semverLatestInRange(version, versionList);
}

if(!finalVersion) {
    console.error('Error: Couldn\'t resolve version for ' + name + ' with range ' + version)
}

Dependency.version = finalVersion;
Dependency.url = 'https://registry.npmjs.org/' + 
    name +
    '/-/' +
    name +
    '-' +
    finalVersion +
    '.tgz'

Dependency;
