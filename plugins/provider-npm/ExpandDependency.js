var name = Depedency.name;
var version = Depedency.version;
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

Depedency.version = finalVersion;

var path = 
    getDepedency('npm', name, finalVersion, 'https://registry.npmjs.org/' + 
        name +
        '/-/' +
        name +
        '-' +
        finalVersion +
        '.tgz');

Depedency;