package html

var (
	SignInPage string = `<!DOCTYPE html>
<html>
  <head>
    <title>Registration page</title>
    <link rel="stylesheet" href="styles.css" />
    <style>
    html {
        height: 100%;
    }
    body {
        display: grid;
        justify-items: center;
        height: 99%;
    }
    .center {
        display: grid;
        margin: auto;
        justify-items: center;
        /*border: solid red 2px;*/
    }
    </style>
  </head>
  <body>
    <div class="center">
      <form class="center" action="{{.SignInRequest}}">
          <label for="login">Login</label>
          <input type="text" id="login" name="login" placeholder="login">
          
          <label for="password">Password</label>
          <input type="text" id="password" name="password" placeholder="*****">
          
          <input type="submit" value="Sign in"><br>
      </form>
      <form action="{{.SignUpPageRequest}}">
          <input type="submit" value="Sign up">
      </form>
    </div>
  </body>
</html>`
	ErrorPage string = `<!DOCTYPE html>
<html>
  <head>
    <title>Registration page</title>
    <link rel="stylesheet" href="styles.css" />
    <style>
    html {
        height: 100%;
    }
    body {
        height: 99%;
    }
    .center_content {
        display: grid;
        margin: auto;
        justify-items: center;
    }
    .form_grid {
        grid-template-columns: repeat(2, 1fr);
    }
    </style>
  </head>
  <body>
    <header class="center_content">
      <h1>{{.Message}}</h1>
    </header>
    <hr>
    <main>
        <div class="center_content">
            <p><b>Error desciption</b></p>
            <form action="{{.SignInPageRequest}}">
              <input type="submit" value="Return">
            </form>
        </div>
    </main>
  </body>
</html>`
	SignUpPage string = `<!DOCTYPE html>
<html>
  <head>
    <title>Registration page</title>
    <link rel="stylesheet" href="styles.css" />
    <style>
    html {
        height: 100%;
    }
    body {
        display: grid;
        justify-items: center;
        height: 99%;
    }
    .center {
        display: grid;
        margin: auto;
        justify-items: center;
    }
    .form_grid {
        grid-template-columns: repeat(2, 1fr);
    }
    .grid_full_line {
        grid-column-start: 1;
        grid-column-end: 3;
    }
    .fird_div_full_line {
        width: 100%;
        display: grid;
        justify-items: center;
    }
    </style>
  </head>
  <body>
    <div class="center">
        <p>SIGN UP</p>
      <form class="center form_grid" action="{{.SignUpRequest}}">
          
        <label for="login">Login</label>
        <input type="text" id="login" name="login" placeholder="login">
        <label for="password">Password</label>
        <input type="text" id="password" name="password" placeholder="*****">
        <div class="grid_full_line">
            <br>
        </div>
        
        <label for="surname">Surname</label>
        <input type="text" id="surname" name="surname" placeholder="Surname">
        <label for="name">Name</label>
        <input type="text" id="name" name="name" placeholder="Name">
        <label for="patronimic">Patronimic</label>
        <input type="text" id="patronimic" name="patronimic" placeholder="Patronimic">
        <label for="passport_code">Passport code</label>
        <input type="text" id="passport_code" name="passport_code" placeholder="0000-00000">
        <label for="birth_date">Birth date</label>
        <input type="text" id="birth_date" name="birth_date" placeholder="01.02.2001">
        <label for="birth_location">Birth location</label>
        <input type="text" id="birth_location" name="birth_location" placeholder="Moscow">
        <label for="pick_up_point">Pick up point</label>
        <input type="text" id="pick_up_point" name="pick_up_point" placeholder="Pick up point name">
        <label for="authority">Authority</label>
        <input type="text" id="authority" name="authority" placeholder="123-321">
        <label for="authority_date">Authority date</label>
        <input type="text" id="authority_date" name="authority_date" placeholder="01.02.2001">
        <label for="registration_adress">Registration address</label>
        <input type="text" id="registration_adress" name="registration_adress" placeholder="Current registration adress">
          
        <label for="inn">INN</label>
        <input type="text" id="inn" name="inn" placeholder="INN">
        <br>
          
        <div class="fird_div_full_line grid_full_line">
            <input type="submit" value="Register" style="width: 100%;">
        </div>
          
      </form>
    </div>
  </body>
</html>`
)
