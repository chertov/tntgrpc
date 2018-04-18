const grpc = require('grpc');
process.exitCode = 1;

const tests = grpc.load(__dirname + '/tests.proto').tests;
const testService = new tests.TestService('localhost:50051', grpc.credentials.createInsecure());
const promises = [];
for (let i = 0; i < 1; i++) {
    promises.push(new Promise((resolve, reject) => {
        testService.NilReturn({}, (err, reply) => {
            if (err) {
                const grpcerr = JSON.parse(err.details);
                if (grpcerr.type === 'MessageNotLuaTable') { resolve(); return; }
            }
            reject('testService.NilReturn');
        })
    }));
    promises.push(new Promise((resolve, reject) => {
        testService.MessageNotLuaTable({}, (err, reply) => {
            if (err) {
                const grpcerr = JSON.parse(err.details);
                if (grpcerr.type === 'MessageNotLuaTable') { resolve(); return; }
            }
            reject('testService.MessageNotLuaTable');
        })
    }));
    promises.push(new Promise((resolve, reject) => {
        testService.NonExistingField({}, (err, reply) => {
            if (err) {
                const grpcerr = JSON.parse(err.details);
                if (grpcerr.type === 'MessageNotExistingField' ) { resolve(); return; }
            }
            reject('testService.NonExistingField');
        })
    }));
    promises.push(new Promise((resolve, reject) => {
        testService.IncorrectFieldType({}, (err, reply) => {
            if (err) {
                const grpcerr = JSON.parse(err.details);
                if (grpcerr.type === 'IncorrectType' && grpcerr.stack[1].type === 'MessageFieldParseFailure' ) { resolve(); return; }
            }
            reject('testService.IncorrectFieldType');
        })
    }));
}

Promise.all(promises)
    .then(() => process.exit(0))
    .catch((ex) => {
        console.log(ex, " failed");
    });
