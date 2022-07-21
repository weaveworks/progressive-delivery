

# Test plan
This document details the test cases for the Progressive Delivery feature of WGE. It is meant to capture:
- the scope of the functionality being tested before a release
- the specific scenarios that need to be tested

## Audience
- Engineers to ensure that the following scenarios are working as described
- Product managers to ensure that the following scenarios are capturing the main functionality that is being delivered

## Scope

The user requirements of this feature are captured in [Notion](https://www.notion.so/weaveworks/Progressive-Delivery-e9dfff425b054b019e72d32a74f553c9#a119b9f5ae5f415783cad691ba3412ed) and [Github](https://github.com/weaveworks/weave-gitops-enterprise/issues/842). The design of the UI screens is captured in [Figma](https://www.figma.com/file/IVHnM9iyeFWpd11evtY8ux/Weave-GitOps?node-id=8620%3A40726). 

Not all of the user requirements will be delivered, see below for details.

## Scenarios

### Flagger not installed ([requirement](https://www.notion.so/weaveworks/Progressive-Delivery-e9dfff425b054b019e72d32a74f553c9#f9106fda720b434883a0e79db3903f62) / [design](https://www.figma.com/file/IVHnM9iyeFWpd11evtY8ux/Weave-GitOps?node-id=8701%3A43258))

1.  Given Flagger is not installed on the management cluster

    When I navigate to the Delivery section of the WGE UI
      
    Then I should see an [onboarding message](https://www.figma.com/file/IVHnM9iyeFWpd11evtY8ux/Weave-GitOps?node-id=8701%3A43258) pointing me to the Flagger docs

### Canary list ([requirement](https://www.notion.so/weaveworks/Progressive-Delivery-e9dfff425b054b019e72d32a74f553c9#f61d40c4b96348798552c1bf8d99440e) / [design](https://www.figma.com/file/IVHnM9iyeFWpd11evtY8ux/Weave-GitOps?node-id=9424%3A51108))

1.  Given Flagger is installed on the management cluster

    And there are no canaries installed on the management cluster

    When I navigate to the Delivery section of the WGE UI

    Then I should see the message `No data to display`.

1.  Given Flagger is installed on the management cluster

    And there is a canary resource installed on the management cluster

    When I navigate to the Delivery section of the WGE UI

    Then I should see a list with a single row representing the canary resource that has been installed

    And the list row should include values for the following:
    - the name of the canary
    - the canary strategy as an icon
    - the status of the canary as an icon and text
    - the cluster that the canary is installed 
    - the namespace of the canary
    - the target of the canary
    - the most recent message of the canary
    - the image that is currently available (including the image registry and the image version)
    - when it was last updated

    And the values should match the state of the canary installed

    And the breadcrump section at the top of the page should indicate that there is one canary in the list


1.   Given Flagger is installed on a leaf cluster

    And there is a canary resource installed on the leaf cluster

    When I navigate to the Delivery section of the WGE UI in the management cluster

    Then I should see a list with a single row representing the canary resource that has been installed

    And the list row should include values for the following:
    - the name of the canary
    - the canary strategy as an icon
    - the status of the canary as an icon and text
    - the cluster that the canary is installed 
    - the namespace of the canary
    - the target of the canary
    - the most recent message of the canary
    - the image that is currently available (including the image registry and the image version)
    - when it was last updated

    And the values should match the state of the canary installed

    And the breadcrump section at the top of the page should indicate that there is one canary in the list

### Canary details - Details tab ([requirement](https://www.notion.so/weaveworks/Progressive-Delivery-e9dfff425b054b019e72d32a74f553c9#f61d40c4b96348798552c1bf8d99440e) / [design](https://www.figma.com/file/IVHnM9iyeFWpd11evtY8ux/Weave-GitOps?node-id=8657%3A42894))

1.  Given Flagger is installed on the management cluster

    And there is a canary resource installed on the management cluster 
    
    And the canary deployment is installed via GitOps

    When I navigate to the Delivery section of the WGE UI

    And click on the name of the canary

    Then I should see the details tab with the following fields:
    - the cluster that the canary is installed 
    - the namespace of the canary
    - the target of the canary
    - the application used to reconcile the target of the canary
    - the canary provider

    And the details tab should also include a status table that shows the current status of the canary along with a nested table that shows the status conditions of the canary

1.  Given Flagger is installed on the management cluster

    And there is a canary resource installed on the management cluster 
    
    And the canary deployment is installed via `kubectl apply`

    When I navigate to the Delivery section of the WGE UI

    And click on the name of the canary

    Then I should see the details tab with the following fields:
    - the cluster that the canary is installed 
    - the namespace of the canary
    - the target of the canary
    - the canary provider

    And the details tab should also include a status table that shows the current status of the canary along with a nested table that shows the status conditions of the canary

### Canary details - Objects tab

1.  Given Flagger is installed on the management cluster

    And there is a canary resource installed on the management cluster

    When I navigate to the Delivery section of the WGE UI

    And click on the name of the canary

    And then click on the Objects tab

    Then I should see a list of all the objects generated by Flagger

    And the list should include at least the following object types:
    - v1/Deployment
    - v1/Service

### Canary details - Events tab ([requirement](https://www.notion.so/weaveworks/Progressive-Delivery-e9dfff425b054b019e72d32a74f553c9#f02709a5fc1c467c910d157de44e296a) / [design](https://www.figma.com/file/IVHnM9iyeFWpd11evtY8ux/Weave-GitOps?node-id=8657%3A42021))


1.  Given Flagger is installed on the management cluster

    And there is a canary resource installed on the management cluster 

    And there are events related to the canary in the last 1 hour

    When I navigate to the Delivery section of the WGE UI

    And click on the name of the canary

    And then click on the Events tab

    Then I should see a list of all events related to the canary that happened over the last 1 hour

    And I should see the following columns for each event:
    - the reason for the event
    - the message of the event
    - the origin of the event
    - the relative time of when the event took place


1.  Given Flagger is installed on the management cluster

    And there is a canary resource installed on the management cluster 

    And there are no events related to the canary in the last 1 hour

    When I navigate to the Delivery section of the WGE UI

    And click on the name of the canary

    And then click on the Events tab

    Then I should see a message `No data to display`

### Canary details - Analysis tab ([requirement](https://github.com/weaveworks/weave-gitops-enterprise/issues/842))

1.  Given Flagger is installed on the management cluster

    And there is a canary resource installed on the management cluster 

    And the canary is configured to use built in metrics for analysis

    When I navigate to the Delivery section of the WGE UI

    And click on the name of the canary

    And then click on the Analysis tab

    Then I should see a list of the metrics that the canary is configured with

    And I should see the following columns for each metric:
    - the metric name
    - the metric template which will be `-` since it uses built in metrics
    - the threshold min value
    - the threshold max value
    - the interval

1.  Given Flagger is installed on the management cluster

    And there is a canary resource installed on the management cluster 

    And the canary is configured to use a custom metric for analysis

    When I navigate to the Delivery section of the WGE UI

    And click on the name of the canary

    And then click on the Analysis tab

    Then I should see a list of the metrics that the canary is configured with

    And I should see the following columns for each metric:
    - the metric name
    - the metric template which is a link that opens the view to the custom metric YAML definition
    - the threshold min value
    - the threshold max value
    - the interval

### Canary details - YAML tab ([requirement](https://www.notion.so/weaveworks/Progressive-Delivery-e9dfff425b054b019e72d32a74f553c9#b9fba13543d448a984d25a44f34f14b5) / [design](https://www.figma.com/file/IVHnM9iyeFWpd11evtY8ux/Weave-GitOps?node-id=8657%3A42392))

1.  Given Flagger is installed on the management cluster

    And there is a canary resource installed on the management cluster 

    When I navigate to the Delivery section of the WGE UI

    And click on the name of the canary

    And then click on the YAML tab

    Then I should see the YAML representation of the canary in numbered lines

### Feature flag ([requirement](https://www.notion.so/weaveworks/Progressive-Delivery-e9dfff425b054b019e72d32a74f553c9#21d6044d867a4624a1a116b83bf260fb))

This requirement is no longer in scope.

### Manual gating docs ([requirement](https://www.notion.so/weaveworks/Progressive-Delivery-e9dfff425b054b019e72d32a74f553c9#05493fd92cb1422086214ad7d8256da0))

Non functional requirement. A [guide](http://docs.gitops.weave.works/docs/next/guides/flagger-manual-gating/) should be available in the Weave GitOps docs that describes how to setup manual gating with Flagger.

### Add Flagger to Flux runtime view ([requirement](https://www.notion.so/weaveworks/Progressive-Delivery-e9dfff425b054b019e72d32a74f553c9#f45ce1a79381401f94fa84c2b9304b91))

Not currently in scope as it requires changes to the Flagger codebase.There is a follow-up [issue](https://github.com/weaveworks/weave-gitops-enterprise/issues/1110) to track this.

### RBAC ([requirement](https://www.notion.so/weaveworks/Progressive-Delivery-e9dfff425b054b019e72d32a74f553c9#aa588eecb393430a898813f54bf74120))

1.  Given Flagger is installed on the management cluster

    And there is a canary resource installed on the management cluster

    And the user account I am logged in as does not have permissions to query for canary resources

    When I navigate to the Delivery section of the WGE UI

    Then I should see the message `No data to display`.


### Realtime data ([requirement](https://www.notion.so/weaveworks/Progressive-Delivery-e9dfff425b054b019e72d32a74f553c9#9e9e4627a01842c8a1a9643e6da062a3))

1.  Given Flagger is installed on the management cluster

    And there is a canary resource installed on the management cluster 

    When I navigate to the Delivery section of the WGE UI

    And I start a progressive rollout of the canary

    Then I should be able to see the status of the canary changing every ~10 seconds.