import { check, sleep, group } from 'k6';
import { GlobalVariables, UUIDs } from './api.v1.global_variables.js';
import { baseURL } from './flag.js';


const glovar = new GlobalVariables(baseURL);
const uuids = new UUIDs;

export function setup() {
    let result = glovar.Find();

    check(result, { 'setup: status is 200': (r) => r.status == 200 });

    JSON.parse(result.body).forEach(e => {
        console.log(`setup: uuid=${e.uuid} name=${e.name} value=${e.value}`)
    })
    return { glovars: JSON.parse(result.body) }
}


export function teardown(data) {
    data.glovars.forEach(e => {
        console.log(`teardown: uuid=${e.uuid} name=${e.name} value=${e.value}`)
        let result = glovar.Update(e.uuid, e.value);
    
        check(result, { 'teardown: status is 200': (r) => r.status == 200 });
    })
}
  

export default function (data) {
    data.glovars.forEach(e => {
        switch (e.uuid) {
            case uuids.ClusterTokenExpirationTime():
                group(`ClusterTokenExpirationTime`, () => {
                    const value_default = e.value
                    const value = '87600h'

                    let result = glovar.Update(e.uuid, value);
                    console.log(`updated body=${result.body}`)
        
                    check(result, { 'check status is 200': (r) => r.status == 200 });
                    check(JSON.parse(result.body), { 'check value': (r) => r.value == value });
        
                    let updated = JSON.parse(result.body)
        
                    group(`check updated`, () => {
                        let result = glovar.Get(e.uuid);
                        console.log(`get body=${result.body}`)
        
                        check(result, { 'status is 200': (r) => r.status == 200 });
                        check(JSON.parse(result.body), { 'check value': (r) => r.value == updated.value });
                    })
                })
                break;
            case uuids.ClientSessionSignatureSecret():
                group(`ClientSessionSignatureSecret`, () => {
                    const value_default = e.value
                    const value = 'YnJvd24gZm94IGp1bXBzIG92ZXIgdGhlIGxhenkgZG9nCg=='

                    let result = glovar.Update(e.uuid, value);
                    console.log(`updated body=${result.body}`)
        
                    check(result, { 'check status is 200': (r) => r.status == 200 });
                    check(JSON.parse(result.body), { 'check value': (r) => r.value == value });
        
                    let updated = JSON.parse(result.body)
        
                    group(`check updated`, () => {
                        let result = glovar.Get(e.uuid);
                        console.log(`get body=${result.body}`)
        
                        check(result, { 'status is 200': (r) => r.status == 200 });
                        check(JSON.parse(result.body), { 'check value': (r) => r.value == updated.value });
                    })
                })
                break;
            case uuids.ClientSessionExpirationTime():
                group(`ClientSessionExpirationTime`, () => {
                    const value_default = e.value
                    const value = '10m'

                    let result = glovar.Update(e.uuid, value);
                    console.log(`updated body=${result.body}`)
        
                    check(result, { 'check status is 200': (r) => r.status == 200 });
                    check(JSON.parse(result.body), { 'check value': (r) => r.value == value });
        
                    let updated = JSON.parse(result.body)
        
                    group(`check updated`, () => {
                        let result = glovar.Get(e.uuid);
                        console.log(`get body=${result.body}`)
        
                        check(result, { 'status is 200': (r) => r.status == 200 });
                        check(JSON.parse(result.body), { 'check value': (r) => r.value == updated.value });
                    })
                })
                break;
            case uuids.ClientConfigServiceValidityPeriod():
                group(`ClientConfigServiceValidityPeriod`, () => {
                    const value_default = e.value
                    const value = '100m'

                    let result = glovar.Update(e.uuid, value);
                    console.log(`updated body=${result.body}`)
        
                    check(result, { 'check status is 200': (r) => r.status == 200 });
                    check(JSON.parse(result.body), { 'check value': (r) => r.value == value });
        
                    let updated = JSON.parse(result.body)
        
                    group(`check updated`, () => {
                        let result = glovar.Get(e.uuid);
                        console.log(`get body=${result.body}`)
        
                        check(result, { 'status is 200': (r) => r.status == 200 });
                        check(JSON.parse(result.body), { 'check value': (r) => r.value == updated.value });
                    })
                })
                break;
        }
    })

    data.glovars.forEach(e => {
        console.log(`restore: uuid=${e.uuid} name=${e.name} value=${e.value}`)
        let result = glovar.Update(e.uuid, e.value);
    
        check(result, { 'restore: status is 200': (r) => r.status == 200 });
    })
}
