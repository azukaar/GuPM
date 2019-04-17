var packageJson = readJsonFile(Dependency.path + '/package.json');
var dependencies = packageJson.dependencies;

Dependency.dependencies = [];

for(depName in dependencies) {
    var depVersion = dependencies[depName];
    Dependency.dependencies.push({
        provider: 'npm',
        name: depName,
        version: depVersion.replace(/(\d) ([\>\<\=\^\~\!])/g, '$1, $2')
    })
}

Dependency;
