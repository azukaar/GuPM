var name = Dependency.name;
var version = Dependency.version;
var finalVersion;

var payload = httpGet('https://registry.npmjs.org/'+name);
var versionList = Object.keys(payload.versions);

// test tags
if(version.match(/^\d*_*\w+[\d\w_]*$/)) {
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

Dependency.version = finalVersion;

var path = 
    getDependency('npm', name, finalVersion, 'https://registry.npmjs.org/' + 
        name +
        '/-/' +
        name +
        '-' +
        finalVersion +
        '.tgz');

// TODO : read package.json

Dependency;