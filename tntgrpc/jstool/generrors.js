const yaml = require('js-yaml');
const fs   = require('fs');
const scriptPath = require('path').dirname(require.main.filename);
let outputDirectory = './../gen/';

process.argv.forEach( val => {
    if (val.startsWith('--outputdir=')) outputDirectory = val.replace('--outputdir=', '');
})
console.log('scriptPath: ', scriptPath);
console.log('outputDirectory: ', outputDirectory);

try {
    let filestr = fs.readFileSync(scriptPath + '/errors.yml', 'utf8');
    filestr = filestr.replace(/(^[ ]+)(ru: )/gmi, '$1$2|\n$1    ');
    filestr = filestr.replace(/(^[ ]+)(en: )/gmi, '$1$2|\n$1    ');
    // fs.writeFileSync(scriptPath + '/errors_fix.yml',filestr);
    const doc = yaml.safeLoad(filestr);
    //fs.writeFileSync('errors.json', JSON.stringify(doc));
    let hpp = "";
    let code = doc.start_code;
    Object.keys(doc.errors).map( (name) => {
        const error = doc.errors[name];
        if (!error.code) error.code = code; code += 1;
        hpp += `class ${name} : public GRPCException {\n`
        { // vars
            const args = error.params.map(arg => `    const std::string ${arg};\n`);
            hpp += `${args.join('')}`;
        }
        const Public = [];
        const Private = [];
        { // string constructor
            let cargs = error.params.map((param) => `const std::string &${param}_`);
            const init_args = error.params.map((param) => `${param}(${param}_)`);

            let str = `    ${name}(${cargs.join(', ')}) : ${init_args.join(', ')} {\n`;
            if (error.lang) Object.keys(error.lang).forEach((langid) => {
                const lang = error.lang[langid];
                // str += `        { std::stringstream ss; ss << ${lang.trim()}; ${langid} = ss.str(); }\n`;
                str += `        { ${langid} = ${lang.trim()}; }\n`;
            });
            str += `    }`;
            Private.push(str);

            cargs = error.params.map((param) => `const std::string &${param}`);
            str =  `   static GRPCError New(${cargs.join(', ')}) {\n`;
            str += `        return std::make_shared<${name}>(${name}(${error.params.join(',')}));\n`;
            str += `   }`;
            Public.push(str);
        }
        if (error.args) error.args.forEach((args) => {
            let params = error.params.map(param => { return {name: param, type:'std::string'}; });
            params = params.map(param => {
                const arg = args[param.name];
                if (arg) {
                    if (arg[0]) param.type = arg[0];
                    if (arg[1]) param.to_str = arg[1];
                }
                return param;
            })

            let cargs = params.map(param => `const ${param.type} &${param.name}_`);
            let str = `    ${name}(${cargs.join(', ')})`;
            const init_args = params.map(param => param.to_str ? `${param.to_str}(${param.name}_)` : `${param.name}_`);
            str += ` : ${name}(${init_args.join(', ')}) {}`;
            Private.push(str);

            cargs = params.map((param) => `const ${param.type} &${param.name}`);
            str = `    static GRPCError New(${cargs.join(', ')}) { return std::make_shared<${name}>(${name}(${error.params.join(',')})); }`;
            Public.push(str);
        });
        hpp += Private.join('\n') + '\n';
        hpp += `public:\n`;
        hpp += Public.join('\n') + '\n';

        hpp += `    virtual int code() const override { return ${error.code}; };\n`;
        hpp += `    virtual std::string type() const override { return "${name}"; };\n`;
        hpp += `    virtual nlohmann::json to_json() const override {\n`;
        const args = error.params.map(param => `{"${param}",${param}}`);
        hpp += `        return {{"errorCode",code()},{"type",type()},{"log",log()},{"loc",location},{"params",{${args.join(',')}}}};\n`;
        hpp += `    }\n`;
        hpp += `};\n`;
    })
    fs.writeFileSync(outputDirectory + '/errors.gen.hpp', hpp);
} catch (e) {
    console.log(e);
    process.exit(1);
}
