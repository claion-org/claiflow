import { check, sleep, group } from 'k6';
import { Webhook } from './api.v1.webhook.js';
import { baseURL } from './flag.js';
import { deepCompare } from './util.js';


const webhook = new Webhook(baseURL);

export function setup() {

}

export function teardown(data) {

}

export default function (data) {
    group(`API TEST WEBHOOK`, () => {
        const uuid               = null;
        const name               = `API-TEST-WEBHOOK-${Date.now()}`;
        const summary            = `API TEST WEBHOOK`;
        const url                = `http:localhost`;
        const method             = `POST`;
        const headers            = null;
        const timeout            = `0s`;
        const conditionValidator = `none`;
        const conditionFilter    = null;
 
        var newWebhook = null
        group(`CreateWebhook`, () => {
            let result = webhook.Create(uuid, name, summary, url, method, headers, timeout, conditionValidator, conditionFilter);
            console.log(`Created body=${result.body}`)
            console.log(`check headers expected=${headers} actual=${headers}`)

            check(result, { 'check status is 200': (r) => r.status == 200 });
            check(JSON.parse(result.body), { 'check name': (r) => r.name == name });
            check(JSON.parse(result.body), { 'check summary': (r) => r.summary == summary });
            check(JSON.parse(result.body), { 'check url': (r) => r.url == url });
            check(JSON.parse(result.body), { 'check method': (r) => r.method == method });
            check(JSON.parse(result.body), { 'check headers': (r) => r.headers == headers });
            check(JSON.parse(result.body), { 'check timeout': (r) => r.timeout == timeout });
            check(JSON.parse(result.body), { 'check conditionValidator': (r) => r.conditionValidator == conditionValidator });
            check(JSON.parse(result.body), { 'check conditionFilter': (r) => r.conditionFilter == conditionFilter });

            newWebhook = JSON.parse(result.body)

            group(`GetWebhook`, () => {
                let result = webhook.Get(newWebhook.uuid, name, summary);
                console.log(`Got body=${result.body}`)

                check(result, { 'check status is 200': (r) => r.status == 200 });
                check(JSON.parse(result.body), { 'check name': (r) => r.name == name });
                check(JSON.parse(result.body), { 'check summary': (r) => r.summary == summary });
                check(JSON.parse(result.body), { 'check url': (r) => r.url == url });
                check(JSON.parse(result.body), { 'check method': (r) => r.method == method });
                check(JSON.parse(result.body), { 'check headers': (r) => r.headers == headers });
                check(JSON.parse(result.body), { 'check timeout': (r) => r.timeout == timeout });
                check(JSON.parse(result.body), { 'check conditionValidator': (r) => r.conditionValidator == conditionValidator });
                check(JSON.parse(result.body), { 'check conditionFilter': (r) => r.conditionFilter == conditionFilter });
            })
        })

        const name_chaned               = `${name}-changed`;
        const summary_chaned            = `${summary}-changed`;
        const url_changed               = `${url}/changed`;
        const method_changed            = `PUT`;
        const headers_changed           = Object.fromEntries(
            new Map( [["Content-Type", ["application/json"]]] ));
        const timeout_changed   = `10s`;
        const conditionValidator_changed = `jq`;
        const conditionFilter_changed = `.status == true`;

        group(`UpdateWebhook (name)`, () => {
            let result = webhook.Update(newWebhook.uuid,  name_chaned, null, null, null, null, null, null, null);
            console.log(`UpdateWebhook (name) body=${result.body}`)


            check(result, { 'status is 200': (r) => r.status == 200 });
            check(JSON.parse(result.body), { 'check name': (r) => r.name == name_chaned });
            check(JSON.parse(result.body), { 'check summary': (r) => r.summary == summary });
            check(JSON.parse(result.body), { 'check url': (r) => r.url == url });
            check(JSON.parse(result.body), { 'check method': (r) => r.method == method });
            check(JSON.parse(result.body), { 'check headers': (r) => r.headers == headers });
            check(JSON.parse(result.body), { 'check timeout': (r) => r.timeout == timeout });
            check(JSON.parse(result.body), { 'check conditionValidator': (r) => r.conditionValidator == conditionValidator });
            check(JSON.parse(result.body), { 'check conditionFilter': (r) => r.conditionFilter == conditionFilter });
      
      
            group(`GetWebhook`, () => {
                let result = webhook.Get(newWebhook.uuid);
                console.log(`Got body=${result.body}`)

                check(result, { 'status is 200': (r) => r.status == 200 });
                check(JSON.parse(result.body), { 'check name': (r) => r.name == name_chaned });
                check(JSON.parse(result.body), { 'check summary': (r) => r.summary == summary });
                check(JSON.parse(result.body), { 'check url': (r) => r.url == url });
                check(JSON.parse(result.body), { 'check method': (r) => r.method == method });
                check(JSON.parse(result.body), { 'check headers': (r) => r.headers == headers });
                check(JSON.parse(result.body), { 'check timeout': (r) => r.timeout == timeout });
                check(JSON.parse(result.body), { 'check conditionValidator': (r) => r.conditionValidator == conditionValidator });
                check(JSON.parse(result.body), { 'check conditionFilter': (r) => r.conditionFilter == conditionFilter });
            })
        })

        group(`UpdateWebhook (summary)`, () => {
            let result = webhook.Update(newWebhook.uuid, null, summary_chaned, null, null, null, null, null, null);
            console.log(`UpdateWebhook (summary) body=${result.body}`)

            check(result, { 'status is 200': (r) => r.status == 200 });
            check(JSON.parse(result.body), { 'check name': (r) => r.name == name_chaned });
            check(JSON.parse(result.body), { 'check summary': (r) => r.summary == summary_chaned });
            check(JSON.parse(result.body), { 'check url': (r) => r.url == url });
            check(JSON.parse(result.body), { 'check method': (r) => r.method == method });
            check(JSON.parse(result.body), { 'check headers': (r) => r.headers == headers });
            check(JSON.parse(result.body), { 'check timeout': (r) => r.timeout == timeout });
            check(JSON.parse(result.body), { 'check conditionValidator': (r) => r.conditionValidator == conditionValidator });
            check(JSON.parse(result.body), { 'check conditionFilter': (r) => r.conditionFilter == conditionFilter });
      
            group(`GetWebhook`, () => {
                let result = webhook.Get(newWebhook.uuid);
                console.log(`Got body=${result.body}`)

                check(result, { 'status is 200': (r) => r.status == 200 });
                check(JSON.parse(result.body), { 'check name': (r) => r.name == name_chaned });
                check(JSON.parse(result.body), { 'check summary': (r) => r.summary == summary_chaned });
                check(JSON.parse(result.body), { 'check url': (r) => r.url == url });
                check(JSON.parse(result.body), { 'check method': (r) => r.method == method });
                check(JSON.parse(result.body), { 'check headers': (r) => r.headers == headers });
                check(JSON.parse(result.body), { 'check timeout': (r) => r.timeout == timeout });
                check(JSON.parse(result.body), { 'check conditionValidator': (r) => r.conditionValidator == conditionValidator });
                check(JSON.parse(result.body), { 'check conditionFilter': (r) => r.conditionFilter == conditionFilter });
            })
        })

        group(`UpdateWebhook (URL)`, () => {
            let result = webhook.Update(newWebhook.uuid, null, null, url_changed, null, null, null, null, null);
            console.log(`UpdateWebhook (URL) body=${result.body}`)

            check(result, { 'status is 200': (r) => r.status == 200 });
            check(JSON.parse(result.body), { 'check name': (r) => r.name == name_chaned });
            check(JSON.parse(result.body), { 'check summary': (r) => r.summary == summary_chaned });
            check(JSON.parse(result.body), { 'check url': (r) => r.url == url_changed });
            check(JSON.parse(result.body), { 'check method': (r) => r.method == method });
            check(JSON.parse(result.body), { 'check headers': (r) => r.headers == headers });
            check(JSON.parse(result.body), { 'check timeout': (r) => r.timeout == timeout });
            check(JSON.parse(result.body), { 'check conditionValidator': (r) => r.conditionValidator == conditionValidator });
            check(JSON.parse(result.body), { 'check conditionFilter': (r) => r.conditionFilter == conditionFilter });
      
            group(`GetWebhook`, () => {
                let result = webhook.Get(newWebhook.uuid);
                console.log(`Got body=${result.body}`)

                check(result, { 'status is 200': (r) => r.status == 200 });
                check(JSON.parse(result.body), { 'check name': (r) => r.name == name_chaned });
                check(JSON.parse(result.body), { 'check summary': (r) => r.summary == summary_chaned });
                check(JSON.parse(result.body), { 'check url': (r) => r.url == url_changed });
                check(JSON.parse(result.body), { 'check method': (r) => r.method == method });
                check(JSON.parse(result.body), { 'check headers': (r) => r.headers == headers });
                check(JSON.parse(result.body), { 'check timeout': (r) => r.timeout == timeout });
                check(JSON.parse(result.body), { 'check conditionValidator': (r) => r.conditionValidator == conditionValidator });
                check(JSON.parse(result.body), { 'check conditionFilter': (r) => r.conditionFilter == conditionFilter });
            })
        })

        group(`UpdateWebhook (method)`, () => {
            let result = webhook.Update(newWebhook.uuid, null, null, null, method_changed, null, null, null, null);
            console.log(`UpdateWebhook (method) body=${result.body}`)

            check(result, { 'status is 200': (r) => r.status == 200 });
            check(JSON.parse(result.body), { 'check name': (r) => r.name == name_chaned });
            check(JSON.parse(result.body), { 'check summary': (r) => r.summary == summary_chaned });
            check(JSON.parse(result.body), { 'check url': (r) => r.url == url_changed });
            check(JSON.parse(result.body), { 'check method': (r) => r.method == method_changed });
            check(JSON.parse(result.body), { 'check headers': (r) => r.headers == headers });
            check(JSON.parse(result.body), { 'check timeout': (r) => r.timeout == timeout });
            check(JSON.parse(result.body), { 'check conditionValidator': (r) => r.conditionValidator == conditionValidator });
            check(JSON.parse(result.body), { 'check conditionFilter': (r) => r.conditionFilter == conditionFilter });
      
            group(`GetWebhook`, () => {
                let result = webhook.Get(newWebhook.uuid);
                console.log(`Got body=${result.body}`)

                check(result, { 'status is 200': (r) => r.status == 200 });
                check(JSON.parse(result.body), { 'check name': (r) => r.name == name_chaned });
                check(JSON.parse(result.body), { 'check summary': (r) => r.summary == summary_chaned });
                check(JSON.parse(result.body), { 'check url': (r) => r.url == url_changed });
                check(JSON.parse(result.body), { 'check method': (r) => r.method == method_changed });
                check(JSON.parse(result.body), { 'check headers': (r) => r.headers == headers });
                check(JSON.parse(result.body), { 'check timeout': (r) => r.timeout == timeout });
                check(JSON.parse(result.body), { 'check conditionValidator': (r) => r.conditionValidator == conditionValidator });
                check(JSON.parse(result.body), { 'check conditionFilter': (r) => r.conditionFilter == conditionFilter });
            })
        })

        group(`UpdateWebhook (headers)`, () => {
            let result = webhook.Update(newWebhook.uuid, null, null, null, null, headers_changed, null, null, null);
            console.log(`UpdateWebhook (headers) body=${result.body}`)

            check(result, { 'status is 200': (r) => r.status == 200 });
            check(JSON.parse(result.body), { 'check name': (r) => r.name == name_chaned });
            check(JSON.parse(result.body), { 'check summary': (r) => r.summary == summary_chaned });
            check(JSON.parse(result.body), { 'check url': (r) => r.url == url_changed });
            check(JSON.parse(result.body), { 'check method': (r) => r.method == method_changed });
            check(JSON.parse(result.body), { 'check headers': (r) => deepCompare(r.headers, headers_changed) });
            check(JSON.parse(result.body), { 'check timeout': (r) => r.timeout == timeout });
            check(JSON.parse(result.body), { 'check conditionValidator': (r) => r.conditionValidator == conditionValidator });
            check(JSON.parse(result.body), { 'check conditionFilter': (r) => r.conditionFilter == conditionFilter });
      
            group(`GetWebhook`, () => {
                let result = webhook.Get(newWebhook.uuid);
                console.log(`Got body=${result.body}`)

                check(result, { 'status is 200': (r) => r.status == 200 });
                check(JSON.parse(result.body), { 'check name': (r) => r.name == name_chaned });
                check(JSON.parse(result.body), { 'check summary': (r) => r.summary == summary_chaned });
                check(JSON.parse(result.body), { 'check url': (r) => r.url == url_changed });
                check(JSON.parse(result.body), { 'check method': (r) => r.method == method_changed });
                check(JSON.parse(result.body), { 'check headers': (r) => deepCompare(r.headers, headers_changed) });
                check(JSON.parse(result.body), { 'check timeout': (r) => r.timeout == timeout });
                check(JSON.parse(result.body), { 'check conditionValidator': (r) => r.conditionValidator == conditionValidator });
                check(JSON.parse(result.body), { 'check conditionFilter': (r) => r.conditionFilter == conditionFilter });
            })
        })


        group(`UpdateWebhook (timeout)`, () => {
            let result = webhook.Update(newWebhook.uuid, null, null, null, null, null, timeout_changed, null, null);
            console.log(`UpdateWebhook (timeout) body=${result.body}`)

            check(result, { 'status is 200': (r) => r.status == 200 });
            check(JSON.parse(result.body), { 'check name': (r) => r.name == name_chaned });
            check(JSON.parse(result.body), { 'check summary': (r) => r.summary == summary_chaned });
            check(JSON.parse(result.body), { 'check url': (r) => r.url == url_changed });
            check(JSON.parse(result.body), { 'check method': (r) => r.method == method_changed });
            check(JSON.parse(result.body), { 'check headers': (r) => deepCompare(r.headers, headers_changed) });
            check(JSON.parse(result.body), { 'check timeout': (r) => r.timeout == timeout_changed });
            check(JSON.parse(result.body), { 'check conditionValidator': (r) => r.conditionValidator == conditionValidator });
            check(JSON.parse(result.body), { 'check conditionFilter': (r) => r.conditionFilter == conditionFilter });
      
            group(`GetWebhook`, () => {
                let result = webhook.Get(newWebhook.uuid);
                console.log(`Got body=${result.body}`)

                check(result, { 'status is 200': (r) => r.status == 200 });
                check(JSON.parse(result.body), { 'check name': (r) => r.name == name_chaned });
                check(JSON.parse(result.body), { 'check summary': (r) => r.summary == summary_chaned });
                check(JSON.parse(result.body), { 'check url': (r) => r.url == url_changed });
                check(JSON.parse(result.body), { 'check method': (r) => r.method == method_changed });
                check(JSON.parse(result.body), { 'check headers': (r) => deepCompare(r.headers, headers_changed) });
                check(JSON.parse(result.body), { 'check timeout': (r) => r.timeout == timeout_changed });
                check(JSON.parse(result.body), { 'check conditionValidator': (r) => r.conditionValidator == conditionValidator });
                check(JSON.parse(result.body), { 'check conditionFilter': (r) => r.conditionFilter == conditionFilter });
            })
        })

        group(`UpdateWebhook (conditionValidator)`, () => {
            let result = webhook.Update(newWebhook.uuid, null, null, null, null, null, null, conditionValidator_changed, conditionFilter_changed);
            console.log(`UpdateWebhook (conditionValidator) body=${result.body}`)

            check(result, { 'status is 200': (r) => r.status == 200 });
            check(JSON.parse(result.body), { 'check name': (r) => r.name == name_chaned });
            check(JSON.parse(result.body), { 'check summary': (r) => r.summary == summary_chaned });
            check(JSON.parse(result.body), { 'check url': (r) => r.url == url_changed });
            check(JSON.parse(result.body), { 'check method': (r) => r.method == method_changed });
            check(JSON.parse(result.body), { 'check headers': (r) => deepCompare(r.headers, headers_changed) });
            check(JSON.parse(result.body), { 'check timeout': (r) => r.timeout == timeout_changed });
            check(JSON.parse(result.body), { 'check conditionValidator': (r) => r.conditionValidator == conditionValidator_changed });
            check(JSON.parse(result.body), { 'check conditionFilter': (r) => r.conditionFilter == conditionFilter_changed });
      
            group(`GetWebhook`, () => {
                let result = webhook.Get(newWebhook.uuid);
                console.log(`Got body=${result.body}`)

                check(result, { 'status is 200': (r) => r.status == 200 });
                check(JSON.parse(result.body), { 'check name': (r) => r.name == name_chaned });
                check(JSON.parse(result.body), { 'check summary': (r) => r.summary == summary_chaned });
                check(JSON.parse(result.body), { 'check url': (r) => r.url == url_changed });
                check(JSON.parse(result.body), { 'check method': (r) => r.method == method_changed });
                check(JSON.parse(result.body), { 'check headers': (r) => deepCompare(r.headers, headers_changed) });
                check(JSON.parse(result.body), { 'check timeout': (r) => r.timeout == timeout_changed });
                check(JSON.parse(result.body), { 'check conditionValidator': (r) => r.conditionValidator == conditionValidator_changed });
                check(JSON.parse(result.body), { 'check conditionFilter': (r) => r.conditionFilter == conditionFilter_changed });
            })
        })

        group(`Undo`, () => {
            const header_undo = Object.fromEntries(
                new Map( [] ));
            const conditionFilter_undo = ``;

            let result = webhook.Update(newWebhook.uuid,  name, summary, url, method, header_undo, timeout, conditionValidator, conditionFilter_undo);
            console.log(`Undo body=${result.body}`)

            check(result, { 'status is 200': (r) => r.status == 200 });
            check(JSON.parse(result.body), { 'check name': (r) => r.name == name });
            check(JSON.parse(result.body), { 'check summary': (r) => r.summary == summary });
            check(JSON.parse(result.body), { 'check url': (r) => r.url == url });
            check(JSON.parse(result.body), { 'check method': (r) => r.method == method });
            check(JSON.parse(result.body), { 'check headers': (r) => r.headers == headers });
            check(JSON.parse(result.body), { 'check timeout': (r) => r.timeout == timeout });
            check(JSON.parse(result.body), { 'check conditionValidator': (r) => r.conditionValidator == conditionValidator });
            check(JSON.parse(result.body), { 'check conditionFilter': (r) => r.conditionFilter == conditionFilter_undo });
      
      
            group(`GetWebhook`, () => {
                let result = webhook.Get(newWebhook.uuid);
                console.log(`Got body=${result.body}`)

                check(result, { 'status is 200': (r) => r.status == 200 });
                check(JSON.parse(result.body), { 'check name': (r) => r.name == name });
                check(JSON.parse(result.body), { 'check summary': (r) => r.summary == summary });
                check(JSON.parse(result.body), { 'check url': (r) => r.url == url });
                check(JSON.parse(result.body), { 'check method': (r) => r.method == method });
                check(JSON.parse(result.body), { 'check headers': (r) => r.headers == headers });
                check(JSON.parse(result.body), { 'check timeout': (r) => r.timeout == timeout });
                check(JSON.parse(result.body), { 'check conditionValidator': (r) => r.conditionValidator == conditionValidator });
                check(JSON.parse(result.body), { 'check conditionFilter': (r) => r.conditionFilter == conditionFilter_undo });
            })
        })

        group(`RemoveWebhook`, () => {
            let result = webhook.Delete(newWebhook.uuid);
            console.log(`Removed body=${result.body}`)

            check(result, { 'check status is 200': (r) => r.status == 200 });

            group(`get removed cluster`, () => {
                let result = webhook.Get(newWebhook.uuid);
                console.log(`Got body=${result.body}`)

                check(result, { 'status is 404': (r) => r.status == 404 });
            })
        })
    })
}