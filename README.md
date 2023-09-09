# Introduction
>**T**emplate **Ma**nager for htm**X** (TMaX)

This is simple library to solve the HTMX specific issue regarding `cold access` to pushed url.

## HTMX `hx-push-url` Scenario 

Now let's imagine you have a `index` page that has button to navigate to a `content`. You're using `hx-get` to transition to that `content`. 

Now if the user chooses to refresh the page, you will be back to the `index` without the `content` part of the page. Let's assume you do not want that because it will be a bad user experience. 

You want if user refresh the page, they will be still in the `content` part of the UI. 

HTMX gives us helpful tools such as `hx-push-url` to push our `hx-get` url towards browser history. Now, even if user refresh the page, it will still accessing the `content` part of the view. 

But now you have bigger problem. 

Let's see what's inside `index` page:
```html
<!DOCTYPE html>
<html lang="en">
<head>
    ... 
</head>
<body id="root-body">
    <p>This is Index body</p>
    <button hx-get="/content" hx-target="#root-body" hx-push-url="true">
        Navigate to Content
    </button>
</body>
</html>
```

And then let's see what's inside `content` part:
```html
<p>This is content body</p>
<button hx-get="/" hx-target="#root-body" hx-push-url="true">
    Back to index
</button>
```

We immediately see that `content` does not looks as complete as `index`. It was not a valid HTML, it lacks head, it lacks body, etc. 

Which is understable, as we want `content` as some sort of partials that can be cheaply sent towards the user. But now because we push the `url` towards browser history, we let user access this incomplete partials.

We can solve this too by checking whether the request was triggered by HTMX or not. If it was sent by HTMX, let's send a partials, and if it was not triggered by HTMX, let's send a full page with the partials inside it. 

To see whether a request were triggered by HTMX, we could just check whether `HX-Request` header exist or not. Easy. 

And to solve the page partials, we could make 2 separate `.html` file that represent the partials and the full page. 

But then we need to introduce 2 files, one for the component only, and the other for the composite for a page. This is fine but there must be neater way to do this. 

This library solve that issue by compositing the view via convention. First, it requries root html that represent the "Page" state:

```html
{{define "root"}}
<!DOCTYPE html>
<html lang="en">

<head>
    {{template "headers"}}
    {{template "head" .}}
</head>

<body id="root-body">
    {{template "body" .}}
</body>

</html>
{{end}}
```

Then we can define a component that satisfy the `root` precondition (look for "head" and "body" component):
```html
{{define "head"}}
<title>mwyndham.dev Index</title>
{{end}}

{{define "body"}}
<p>This is Index body</p>
<button hx-get="/content" hx-target="#root-body" hx-push-url="true">
    Navigate to Content
</button>
{{end}}
```

This way, you can call the component state by it's name:
```go
return c.Render(http.StatusOK, "index", nil)
```
or by it's root state:
```go
return c.Render(http.StatusOK, "root#index", nil)
```

## Using the Library

### Labstack Echo

The function itself is look like this:
```go
func NewEchoTemplateRenderer(e *echo.Echo, rootName, viewName, componentPath string, viewPaths ...string)
```
Where each params is:
- `e` -> instance of Echo Server
- `rootName` the name definition of the root layout. For example, if the root layout template is defined as "root" then you should name it "root" too. 
- `viewName` is the name definition of the component. In the example above, the component is defined as "body", so we put the value "body" here. 
- `componentPath` is the directory where we put all of the non-page composable components (including the root.html). It capable to recursively looking for html file so you can arrange the files as you liked as long as they have common parent which you specify here. 
- `viewPaths` is variadic value where you can put one or more path for your page-composable components. 

The full example for those function would be like this:
```go
tmax.NewEchoTemplateRenderer(e, "root", "body", "public/components", "public/views") // Don't forget that the last one is variadic
```



