describe("Kubernetes UI Graph",function(){it("should have all the expected components loaded",function(){browser.get("http://localhost:8000"),expect(browser.getTitle()).toEqual("Kubernetes UI");var e=element(by.id("tab_002"));expect(e).toBeDefined(),e.click(),expect(browser.getLocationAbsUrl()).toBe("/graph/");var t=element(by.id("ExpandCollapse"));expect(t).toBeDefined();var i=element(by.id("ToggleSelect"));expect(i).toBeDefined();var c=element(by.id("ToggleSource"));expect(c).toBeDefined();var o=element(by.id("PollOnce"));expect(o).toBeDefined(),c.click(),o.click();var l=element(by.css("svg"));expect(l).toBeDefined()}),it("should have all the details pane working correctly",function(){browser.get("http://localhost:8000"),expect(browser.getTitle()).toEqual("Kubernetes UI");var e=element(by.id("tab_002"));expect(e).toBeDefined(),e.click(),expect(browser.getLocationAbsUrl()).toBe("/graph/");var t=element(by.id("toggleDetails"));expect(t).toBeDefined(),expect(element(by.repeater("type in getLegendLinkTypes()"))).toBeDefined();var i=element(by.id("details"));expect(i).toBeDefined(),expect(i.isDisplayed()).toBe(!1),t.click(),expect(i.isDisplayed()).toBe(!0)}),it("should have all the graph working correctly",function(){browser.get("http://localhost:8000"),expect(browser.getTitle()).toEqual("Kubernetes UI");var e=element(by.id("tab_002"));expect(e).toBeDefined(),e.click(),expect(browser.getLocationAbsUrl()).toBe("/graph/");var t=element(by.css("d3-visualization svg"));expect(t).toBeDefined(),expect(element(by.css("d3-visualization svg marker"))).toBeDefined(),expect(element(by.css("d3-visualization svg text"))).toBeDefined(),expect(element(by.css("d3-visualization svg path"))).toBeDefined(),expect(element(by.css("d3-visualization svg image"))).toBeDefined(),expect(element(by.css("d3-visualization svg circle"))).toBeDefined();var i=element(by.id("ToggleSource"));expect(i).toBeDefined();var c=element(by.id("PollOnce"));expect(c).toBeDefined(),i.click(),c.click(),by.addLocator("datumIdMatches",function(e,t,i){var c=[];return window.d3.selectAll("circle").each(function(t){t&&t.id===e&&c.push(this)}),c});var o=element(by.datumIdMatches("Pod:redis-slave-controller-vi7hv"));expect(o).toBeDefined(),o.click();var l=element(by.id("details"));expect(l).toBeDefined(),expect(l.isDisplayed()).toBe(!0),expect(element(by.repeater("(tag, value) in getSelectionDetails()"))).toBeDefined();var n=element(by.id("inspectBtn"));expect(n).toBeDefined(),n.click(),expect(browser.getLocationAbsUrl()).toBe("/graph/inspect"),expect(element(by.repeater("(key, val) in json track by $index"))).toBeDefined();var a=element(by.id("backButton"));expect(a).toBeDefined(),a.click(),expect(browser.getLocationAbsUrl()).toBe("/graph/")})});