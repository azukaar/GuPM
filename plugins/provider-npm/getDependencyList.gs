var result = [];

for(depName in PackageConfig.dependencies) {
    var depVersion = PackageConfig.dependencies[depName];
    result.push({
        provider: 'npm',
        name: depName,
        version: depVersion
            .replace(/(\d) ([\>\<\=\^\~\!])/g, '$1, $2')
            .replace(/^\^0/, '~0')
    })
}

for(depName in PackageConfig.devDependencies) {
    var depVersion = PackageConfig.devDependencies[depName];
    result.push({
        provider: 'npm',
        name: depName,
        version: depVersion
            .replace(/(\d) ([\>\<\=\^\~\!])/g, '$1, $2')
            .replace(/^\^0/, '~0')
    })
}


result;