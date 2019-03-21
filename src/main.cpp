#include <iostream>
#include <fstream>
#include <vector>

#include "../deps/rapidjson/document.h"
#include "../deps/rapidjson/writer.h"
#include "../deps/rapidjson/stringbuffer.h"

#include <stdio.h>
#include <unistd.h>
#include <sys/wait.h>

using namespace rapidjson;
using namespace std;

Document getJson(string path) {
		ifstream packagejson;
		string jsonString;

		packagejson.open (path.c_str());
		getline (packagejson, jsonString, (char) packagejson.eof());
		packagejson.close();

    Document d;
    d.Parse(jsonString.c_str());

		return d;
}

void download(string dep, string version) {
		string curl = string("curl -o node_modules/" + dep + ".tgz https://registry.npmjs.org/" + dep + "/" + dep + "-" + version + ".tgz");
		cout << curl << endl;
		system(curl.c_str());
}

int main(int argc, char *argv[]) {
		Document package = getJson("package.json");
		const Value& deps = package["dependencies"];
		assert(deps.IsObject());
		int nbProc = 0;

		for (Value::ConstMemberIterator itr = deps.MemberBegin(); itr != deps.MemberEnd(); ++itr) {
			string dep = string(itr->name.GetString());
			string version = deps[dep.c_str()].GetString();
			
			if(fork() == 0) {
				nbProc++;
				download(dep, version);
				_exit(0);
			}

			pid_t child[nbProc];
			int status[nbProc];
			int i;

			for(i=0;i<nbProc;++i)
				child[i] = wait(&status[i]);

			for(i=0;i<nbProc;++i)
				printf("Exit = %d, child = %d\n", WEXITSTATUS(status[i]), child[i]);
		}

    return 0;
}