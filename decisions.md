## Rest API

Our API guidelines build on the [Zalando RESTful API and Event Guidelines](https://opensource.zalando.com/restful-api-guidelines/).

On top of that the following rules are established.

### 9. REST Basics - HTTP requests

#### POST

##### Change #1

If a POST endpoint aims at creating new resources, then the ID of the new resource MUST NOT be present in the payload.

_Reason:_ This change makes it harder to accidentally re-create resources, and makes clear that the resource identifier management is always under the control of the service and not the client.


#### PUT

##### Change #1

The Zalando guidelines ALLOW not having resource identifiers in the payload, but only in path. AbterGo does not allow that. The payload MUST always contain the ID, and it MUST be checked against the ID received as a path parameter. In case of an ID mismatch, 400 Bad request must be returned.

_Reason:_ This change makes it harder to accidentally update existing resources.

##### Change #2

Also, whereas the Zalando guidelines ALLOW using ETags, AbterGo makes this mandatory. Not complying with this rule MAY break future assumptions.

_Reason:_ This change enables certain features to work regardless of the resources they operate on. (Example middlewares may implement custom logic which require E-Tags to be sent along with PUT requests.)