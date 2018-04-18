const grpc = require('grpc');
process.exitCode = 1;

const tests = grpc.load(__dirname + '/tests.proto').tests;
const doesNotExistsService = new tests.DoesNotExistsService('localhost:50051', grpc.credentials.createInsecure());
const testService = new tests.TestService('localhost:50051', grpc.credentials.createInsecure());

const promises = [];
for (let i = 0; i < 1000; i++) {
    let p;
    p = new Promise((resolve, reject) => {
        doesNotExistsService.Procedure({}, (err, reply) => {
            if (err && err.code === 13) {
                const grpcerr = JSON.parse(err.details);
                if (grpcerr.type === 'ServiceNotFound') { resolve(); return; }
            }
            reject('doesNotExistsService.Procedure failed');
        })
    });
    promises.push(p);
    p = new Promise((resolve, reject) => {
        testService.DoesNotExistsProcedure({}, (err, reply) => {
            if (err && err.code === 13) {
                const grpcerr = JSON.parse(err.details);
                if (grpcerr.type === 'ProcedureNotFound') { resolve(); return; }
            }
            reject('testService.DoesNotExistsProcedure failed');
        })
    });
    promises.push(p);
    p = new Promise((resolve, reject) => {
        testService.ExistsProcedure({message: 'm1'}, (err, reply) => {
            if (!err && reply && reply.message === 'm1') { resolve(); return; }
            reject('testService2.ExistsProcedure failed');
        });
    });
    promises.push(p);
    p = new Promise((resolve, reject) => {
        testService.ErrorWhileCallProcedure({}, (err, reply) => {
            if (err && err.code === 13) {
                const grpcerr = JSON.parse(err.details);
                if (grpcerr.type === 'CallProcedureError') { resolve(); return; }
            }
            reject('testService.ErrorWhileCallProcedure failed');
        })
    });
    promises.push(p);
}

Promise.all(promises)
    .then(() => process.exit(0))
    .catch((ex) => {
        console.log('test failed: ', ex);
    });
