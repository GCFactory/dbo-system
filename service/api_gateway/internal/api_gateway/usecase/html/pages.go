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
      <form class="center" id="signInForm">
          <label for="login">Login</label>
          <input type="text" id="login" name="login" placeholder="login" required>
          
          <label for="password">Password</label>
          <input type="text" id="password" name="password" placeholder="*****" required>
          
			<div id="message"></div>
			<br>
          <input type="submit" value="Sign in" id="signInButton"><br>
      </form>
      <form action="{{.SignUpPageRequest}}">
          <input type="submit" value="Sign up">
      </form>
    </div>

	<script>
		document.getElementById('signInForm').addEventListener('submit', function(event) {
		event.preventDefault(); // Предотвращаем стандартную отправку формы
	
		var submitButton = document.getElementById('signInButton');
		submitButton.disabled = true; // Отключаем кнопку
		submitButton.innerText = 'Check login and password...'; // Меняем текст кнопки
	
		var formData = {
			login: document.querySelector('input[name="login"]').value,
			password: document.querySelector('input[name="password"]').value
		};
	
		// Отправляем POST-запрос с учётом cookie
		fetch('{{.SignInRequest}}', {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json' // Указываем, что отправляем JSON
			},
			body: JSON.stringify(formData), // Преобразуем объект в JSON
			credentials: 'include' // Важно: включаем отправку и получение cookie
		})
		.then(response => {
			if (response.ok) {
				return response.json(); // Парсим JSON-ответ
			} else {
				return response.json().then(errorData => {
					throw new Error(errorData.error || 'Sending post request check login and password error');
				});
			}
		})
		.then(data => {
			if (data.success) {
				window.location.href = '{{.HomePageRequest}}'; // Перенаправляем на GET-страницу
			} else {
				// Если есть ошибка, выводим сообщение
				document.getElementById('message').innerText = data.error;
			}
		})
		.catch(error => {
			document.getElementById('message').innerText = error;
			document.getElementById('message').style.color = 'red'; // Делаем текст красным
			console.error('Error:', error);
		})
		.finally(() => {
			submitButton.disabled = false; // Включаем кнопку обратно
			submitButton.innerText = 'Sign in';
		});
	});
	</script>

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
      <h1>ERROR</h1>
    </header>
    <hr>
    <main>
        <div class="center_content">
            <p><b>{{.Message}}</b></p>
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
      <form class="center form_grid" id="signUpForm">
          
        <label for="login">Login</label>
        <input type="text" id="login" name="login" placeholder="login" required>
        <label for="password">Password</label>
        <input type="text" id="password" name="password" placeholder="*****" required>
        <div class="grid_full_line">
            <br>
        </div>
        
        <label for="surname">Surname</label>
        <input type="text" id="surname" name="surname" placeholder="Surname" required>
        <label for="name">Name</label>
        <input type="text" id="name" name="name" placeholder="Name" required>
        <label for="patronymic">Patronymic</label>
        <input type="text" id="patronymic" name="patronymic" placeholder="Patronymic" required>
        <label for="passport_series">Passport series</label>
        <input type="text" id="passport_series" name="passport_series" placeholder="0000" required>
        <label for="passport_number">Passport series</label>
        <input type="text" id="passport_number" name="passport_number" placeholder="000000" required>
        <label for="birth_date">Birth date</label>
        <input type="text" id="birth_date" name="birth_date" placeholder="01-02-2001" required>
        <label for="birth_location">Birth location</label>
        <input type="text" id="birth_location" name="birth_location" placeholder="Moscow" required>
        <label for="passport_pick_up_point">Pick up point</label>
        <input type="text" id="passport_pick_up_point" name="passport_pick_up_point" placeholder="Pick up point name" required>
        <label for="passport_authority">Authority</label>
        <input type="text" id="passport_authority" name="passport_authority" placeholder="123-321" required>
        <label for="passport_authority_date">Authority date</label>
        <input type="text" id="passport_authority_date" name="passport_authority_date" placeholder="01-02-2001" required>
        <label for="passport_registration_address">Registration address</label>
        <input type="text" id="passport_registration_address" name="passport_registration_address" placeholder="Current registration address" required>
          
        <label for="inn">INN</label>
        <input type="text" id="inn" name="inn" placeholder="01234567890123456789" required>
        <br>
          
		<div id="message"></div>
		<br>
        <div class="fird_div_full_line grid_full_line">
            <input type="submit" value="Register" style="width: 100%;" id="signUpButton">
        </div>
          
      </form>
    </div>

	<script>
		document.getElementById('signUpForm').addEventListener('submit', function(event) {
		event.preventDefault(); // Предотвращаем стандартную отправку формы
	
		var submitButton = document.getElementById('signUpButton');
		submitButton.disabled = true; // Отключаем кнопку
		submitButton.innerText = 'Register user...'; // Меняем текст кнопки
	
		var formData = {
			login: document.querySelector('input[name="login"]').value,
			password: document.querySelector('input[name="password"]').value,
			surname: document.querySelector('input[name="surname"]').value,
			name: document.querySelector('input[name="name"]').value,
			patronymic: document.querySelector('input[name="patronymic"]').value,
			passport_series: document.querySelector('input[name="passport_series"]').value,
			passport_number: document.querySelector('input[name="passport_number"]').value,
			birth_date: document.querySelector('input[name="birth_date"]').value,
			birth_location: document.querySelector('input[name="birth_location"]').value,
			passport_pick_up_point: document.querySelector('input[name="passport_pick_up_point"]').value,
			passport_authority: document.querySelector('input[name="passport_authority"]').value,
			passport_authority_date: document.querySelector('input[name="passport_authority_date"]').value,
			passport_registration_address: document.querySelector('input[name="passport_registration_address"]').value,
			inn: document.querySelector('input[name="inn"]').value
		};
	
		// Отправляем POST-запрос с учётом cookie
		fetch('{{.SignUpRequest}}', {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json' // Указываем, что отправляем JSON
			},
			body: JSON.stringify(formData), // Преобразуем объект в JSON
			credentials: 'include' // Важно: включаем отправку и получение cookie
		})
		.then(response => {
			if (response.ok) {
				return response.json(); // Парсим JSON-ответ
			} else {
				return response.json().then(errorData => {
					throw new Error(errorData.error || 'Sending post request register user error');
				});
			}
		})
		.then(data => {
			if (data.success) {
				window.location.href = '{{.HomePageRequest}}'; // Перенаправляем на GET-страницу
			} else {
				// Если есть ошибка, выводим сообщение
				document.getElementById('message').innerText = data.error;
			}
		})
		.catch(error => {
			document.getElementById('message').innerText = error;
			document.getElementById('message').style.color = 'red'; // Делаем текст красным
			console.error('Error:', error);
		})
		.finally(() => {
			submitButton.disabled = false; // Включаем кнопку обратно
			submitButton.innerText = 'Sign up';
		});
	});
	</script>

  </body>
</html>`
	HomePage string = `<!DOCTYPE html>
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
    header {
        display: grid;
        grid-auto-flow: column; /* Arrange items in a row */
        justify-content: end; /* Align items to the right */
        width: 100%;
    }
    .center_content {
        display: grid;
        margin: auto;
        justify-items: center;
    }
    .pre-tab {
        white-space: pre; /* Сохраняем изначальное форматирование */
    }
    .form_grid {
        grid-template-columns: repeat(2, 1fr);
    }
    .form_grid_4 {
        grid-template-columns: repeat(4, 1fr);
    }
    .bold {
        
    }
    table {
        border: collapse;
        width: 100%;
    }
    thead {
      background-color: rgb(228 240 245);
    }
    td {
        text-align: center;
    }
    tbody tr:nth-child(odd) {
      background-color: rgb(123,123,123);
      color: #fff;
    }
    </style>
  </head>
  <body>
    <header>
      <div>
          <p>{{.Login}}</p>
      </div>
      <p class="pre-tab">   |   </p>
      <div class="center_content">
          <form class="center_content" id="signOutForm">
              <input type="submit" value="Sign out" id="signOutButton">
				<div id="message"></div>
          </form>
      </div>
    </header>
    <hr>
    <main>
        <div class="center_content">
            <p><b>USER DATA</b></p>
        </div>
        <div class="center_content form_grid">
            
            <p><b>Surname</b></p>
            <p>{{.Surname}}</p>
            
            <p><b>Name</b></p>
            <p>{{.Name}}</p>
            
            <p><b>Patronymic</b></p>
            <p>{{.Patronymic}}</p>
            
            <p><b>INN</b></p>
            <p>{{.INN}}</p>
            
            <p><b>Passport code</b></p>
            <p>{{.PassportCode}}</p>
            
            <p><b>Birth date</b></p>
            <p>{{.BirthDate}}</p>
            
            <p><b>Birth location</b></p>
            <p>{{.BirthLocation}}</p>
            
            <p><b>Pick up point</b></p>
            <p>{{.PickUpPoint}}</p>
            
            <p><b>Authority</b></p>
            <p>{{.Authority}}</p>
            
            <p><b>Authority date</b></p>
            <p>{{.AuthorityDate}}</p>
            
            <p><b>Registration address</b></p>
            <p>{{.RegistrationAddress}}</p>
            
        </div>
        <hr>
        <div class="center_content">
            <p><b>ACCOUNTS</b></p>
            <form action="{{.CreateAccountRequest}}">
				<input type="hidden" name="user_id" value="{{.UserId}}">
                <input type="submit" value="Open account">
            </form>
        </div>
        <div class="center_content">
            <table>
                <thead>
                    <tr>
                        <th scope="col">Account name</th>
                        <th scope="col">Status</th>
                        <th scope="col">Account cache</th>
                        <th scope="col">Account operations</th>
                    </tr>
                </thead>
                <tbody>
					{{.ListOfAccounts}}
                </tbody>
            </table>
        </div>
    </main>

	<script>
		document.getElementById('signOutForm').addEventListener('submit', function(event) {
		event.preventDefault(); // Предотвращаем стандартную отправку формы
	
		var submitButton = document.getElementById('signOutButton');
		submitButton.disabled = true; // Отключаем кнопку
		submitButton.innerText = 'Exiting...'; // Меняем текст кнопки
	
		// Отправляем POST-запрос с учётом cookie
		fetch('{{.SignOutRequest}}', {
			method: 'POST',
			credentials: 'include' // Важно: включаем отправку и получение cookie
		})
		.then(response => {
			if (response.ok) {
				return response.json(); // Парсим JSON-ответ
			} else {
				return response.json().then(errorData => {
					throw new Error(errorData.error || 'Sending post request register user error');
				});
			}
		})
		.then(data => {
			if (data.success) {
				window.location.href = '{{.SignInPageRequest}}'; // Перенаправляем на GET-страницу
			} else {
				// Если есть ошибка, выводим сообщение
				document.getElementById('message').innerText = data.error;
			}
		})
		.catch(error => {
			document.getElementById('message').innerText = error;
			document.getElementById('message').style.color = 'red'; // Делаем текст красным
			console.error('Error:', error);
		})
		.finally(() => {
			submitButton.disabled = false; // Включаем кнопку обратно
			submitButton.innerText = 'Sign out';
		});
	});
	</script>

  </body>
</html>`
	AccountOperationPage string = `<!DOCTYPE html>
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
    header {
        display: grid;
        grid-auto-flow: column; /* Arrange items in a row */
        justify-content: end; /* Align items to the right */
        width: 100%;
    }
    .center_content {
        display: grid;
        margin: auto;
        justify-items: center;
    }
    .form_grid {
        grid-template-columns: repeat(2, 1fr);
    }
    .pre-tab {
        white-space: pre; /* Сохраняем изначальное форматирование */
    }
    </style>
  </head>
  <body>
    <header>
      <div>
          <p>{{.Login}}</p>
      </div>
      <p class="pre-tab">   |   </p>
      <div class="center_content">
          <form class="center_content" id="signOutForm">
              <input type="submit" value="Sign out" id="signOutButton">
				<div id="message"></div>
          </form>
      </div>
    </header>
    <hr>
    <main>
        <div class="center_content">
            <h1>{{.OperationName}}</h1>
        </div>

		{{.Operation}}

        <br>        
        <div>
            <form class="center_content" action="{{.ReturnRequest}}">
                <input type="submit" value="Return">
            </form>
        </div>
    </main>

	{{.OperationScript}}

	<script>
		document.getElementById('signOutForm').addEventListener('submit', function(event) {
		event.preventDefault(); // Предотвращаем стандартную отправку формы
	
		var submitButton = document.getElementById('signOutButton');
		submitButton.disabled = true; // Отключаем кнопку
		submitButton.innerText = 'Exiting...'; // Меняем текст кнопки
	
		// Отправляем POST-запрос с учётом cookie
		fetch('{{.SignOutRequest}}', {
			method: 'POST',
			credentials: 'include' // Важно: включаем отправку и получение cookie
		})
		.then(response => {
			if (response.ok) {
				return response.json(); // Парсим JSON-ответ
			} else {
				return response.json().then(errorData => {
					throw new Error(errorData.error || 'Sending post request register user error');
				});
			}
		})
		.then(data => {
			if (data.success) {
				window.location.href = '{{.SignInPageRequest}}'; // Перенаправляем на GET-страницу
			} else {
				// Если есть ошибка, выводим сообщение
				document.getElementById('message').innerText = data.error;
			}
		})
		.catch(error => {
			document.getElementById('message').innerText = error;
			document.getElementById('message').style.color = 'red'; // Делаем текст красным
			console.error('Error:', error);
		})
		.finally(() => {
			submitButton.disabled = false; // Включаем кнопку обратно
			submitButton.innerText = 'Sign out';
		});
	});
	</script>

  </body>
</html>`
	AdminPage string = `<!DOCTYPE html>
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
    header {
        display: grid;
        margin: auto;
        justify-items: center;
        width: 100%;
    }
    .center_content {
        display: grid;
        margin: auto;
        justify-items: center;
    }
    .form_grid {
        grid-template-columns: repeat(2, 1fr);
    }
    .pre-tab {
        white-space: pre; /* Сохраняем изначальное форматирование */
    }
    table {
        border: collapse;
        width: 100%;
    }
    thead {
      background-color: rgb(228 240 245);
    }
    td {
        text-align: center;
    }
    tbody tr:nth-child(odd) {
      background-color: rgb(123,123,123);
      color: #fff;
    }
    </style>
  </head>
  <body>
    <header>
      <h1>ADMIN PAGE</h1>
    </header>
    <hr>
    <main>
      <div class="center_content">
        <form action="{{.GetOperationsRequest}}">
          
          <label for="start">Time begin:</label>
          <input type="text" id="start" name="start" placeholder="01-02-2001 12:12:12">
          
          <label for="end">Time end:</label>
          <input type="text" id="end" name="end" placeholder="01-02-2001 12:12:12">
          
          <input type="submit" value="Get operations">
          
        </form>
      </div>
        <table>
            <thead>
                <tr>
                    <th scope="col">Operation id</th>
                    <th scope="col">Operation name</th>
                    <th scope="col">Operation status</th>
                    <th scope="col">Time start</th>
                    <th scope="col">Time end</th>
                    <th scope="col">Operation tree</th>
                </tr>
            </thead>
            <tbody>
                
				{{.Operations}}

            </tbody>
        </table>
    </main>
  </body>
</html>`
)

var (
	HomePageAccount string = `
		<tr>
			<td>{{.Name}}</td>
			<td>{{.Status}}</td>
			<td>{{.Cache}}</td>
			<td>
				<div class="center_content form_grid_4">
					 <form action="{{.GetCreditsRequest}}">
						<input type="hidden" name="account_id" value="{{.AccountId}}">
						<input type="submit" value="Get credits" {{if .Disabled -}} disabled {{else -}} {{end}}>
					</form>
					<form action="{{.AddCacheRequest}}">
						<input type="hidden" name="account_id" value="{{.AccountId}}">
						<input type="submit" value="Add cache" {{if .Disabled -}} disabled {{else -}} {{end}}>
					</form>
					<form action="{{.ReduceCacheRequest}}">
						<input type="hidden" name="account_id" value="{{.AccountId}}">
						<input type="submit" value="Reduce cache" {{if .Disabled -}} disabled {{else -}} {{end}}>
					</form>
					<form action="{{.CloseAccountRequest}}">
						<input type="hidden" name="account_id" value="{{.AccountId}}">
						<input type="submit" value="Close account" {{if .Disabled -}} disabled {{else -}} {{end}}>
					</form>   
				</div>
				
			</td>
		</tr>
	`
	AccountOperationCreateAccount string = `
        <div>
            <form class="center_content" id="createAccountForm">
                
                <label for="name"><b>Account name</b></label>
                <input type="text" id="name" name="name" placeholder="Name" required>
                
                <label for="culc_number"><b>Culc number</b></label>
                <input type="text" id="culc_number" name="culc_number" placeholder="40705810990123456789" required>
                
                <label for="corr_number"><b>Corr number</b></label>
                <input type="text" id="corr_number" name="corr_number" placeholder="30125810502500000025" required>
                
                <label for="bic"><b>BIC</b></label>
                <input type="text" id="bic" name="bic" placeholder="245025025" required>
                
                <label for="cio"><b>CIO</b></label>
                <input type="text" id="cio" name="cio" placeholder="509910012" required>
				
				<div id="operationMessage"></div>
                
                <input type="submit" value="Create" id="createAccountButton">
                
            </form>
        </div>
	`
	AccountOperationCloseAccount string = `        
		<div>
            <form class="center_content" id="closeForm">
				<input type="hidden" name="account_id" value="{{.AccountId}}">
                <input type="submit" value="Confirm" id="closeButton">
            </form>
        </div>
	`
	AccountOperationGetCredits string = `
		<div class="center_content form_grid">
		
		    <p><b>Account name</b></p>
		    <p>{{.Name}}</p>
		
		    <p><b>Status</b></p>
		    <p>{{.Status}}</p>
		
		    <p><b>Amount</b></p>
		    <p>{{.Amount}}</p>
		
		    <p><b>Culc number</b></p>
		    <p>{{.CulcNumber}}</p>
		
		    <p><b>Corr number</b></p>
		    <p>{{.CorrNumber}}</p>
		
		    <p><b>BIC</b></p>
		    <p>{{.BIC}}</p>
		
		    <p><b>CIO</b></p>
		    <p>{{.CIO}}</p>
		
		</div>
	`
	AccountOperationAddCache string = `
	        <div>
            <form class="center_content" id="formAccountMoney">
				<input type="hidden" name="account_id" value="{{.AccountId}}">
                <lable for="money">Money</lable>
                <input type="text" id="money" name="money">
				<div id="operationMessage"></div>
                <input type="submit" value="Confirm" id="accountMoneyButton">
            </form>    
        </div>`
	AccountOperationWidthCache string = AccountOperationAddCache
	AdminOperation             string = `
	<tr>
		<td>{{.Id}}</td>
		<td>{{.Name}}</td>
		<td>{{.Status}}</td>
		<td>{{.Begin}}</td>
		<td>{{.End}}</td>
		<td>
			<img src="{{.ImagePath}}" 
			style="width:100px; height: 100px;"
			alt="operation_graph">
		</td>
	</tr>
`
)

var (
	ScriptOpenAccountOperation string = `
		<script>
			document.getElementById('createAccountForm').addEventListener('submit', function(event) {
			event.preventDefault(); // Предотвращаем стандартную отправку формы
		
			var submitButton = document.getElementById('createAccountButton');
			submitButton.disabled = true; // Отключаем кнопку
			submitButton.innerText = 'Creating...'; // Меняем текст кнопки
		
			var formData = {
				name: document.querySelector('input[name="name"]').value,
				culc_number: document.querySelector('input[name="culc_number"]').value,
				corr_number: document.querySelector('input[name="corr_number"]').value,
				bic: document.querySelector('input[name="bic"]').value,
				cio: document.querySelector('input[name="cio"]').value,
			};

			// Отправляем POST-запрос с учётом cookie
			fetch('{{.OperationRequest}}', {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json' // Указываем, что отправляем JSON
				},
				body: JSON.stringify(formData), // Преобразуем объект в JSON
				credentials: 'include' // Важно: включаем отправку и получение cookie
			})
			.then(response => {
				if (response.ok) {
					return response.json(); // Парсим JSON-ответ
				} else {
					return response.json().then(errorData => {
						throw new Error(errorData.error);
					});
				}
			})
			.then(data => {
				if (data.success) {
					window.location.href = '{{.HomePageRequest}}'; // Перенаправляем на GET-страницу
				} else {
					// Если есть ошибка, выводим сообщение
					document.getElementById('operationMessage').innerText = data.error;
				}
			})
			.catch(error => {
				document.getElementById('operationMessage').innerText = error;
				document.getElementById('operationMessage').style.color = 'red'; // Делаем текст красным
				console.error('Error:', error);
			})
			.finally(() => {
				submitButton.disabled = false; // Включаем кнопку обратно
				submitButton.innerText = 'Create';
			});
		});
		</script>
`
	ScriptCloseAccountOperation string = `
		<script>
			document.getElementById('closeForm').addEventListener('submit', function(event) {
			event.preventDefault(); // Предотвращаем стандартную отправку формы
		
			var submitButton = document.getElementById('closeButton');
			submitButton.disabled = true; // Отключаем кнопку
			submitButton.innerText = 'Closing...'; // Меняем текст кнопки
		
			var formData = {
				account_id: document.querySelector('input[name="account_id"]').value
			};

			// Отправляем POST-запрос с учётом cookie
			fetch('{{.OperationRequest}}', {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json' // Указываем, что отправляем JSON
				},
				body: JSON.stringify(formData), // Преобразуем объект в JSON
				credentials: 'include' // Важно: включаем отправку и получение cookie
			})
			.then(response => {
				if (response.ok) {
					return response.json(); // Парсим JSON-ответ
				} else {
					return response.json().then(errorData => {
						throw new Error(errorData.error);
					});
				}
			})
			.then(data => {
				if (data.success) {
					window.location.href = '{{.HomePageRequest}}'; // Перенаправляем на GET-страницу
				} else {
					// Если есть ошибка, выводим сообщение
					document.getElementById('operationMessage').innerText = data.error;
				}
			})
			.catch(error => {
				document.getElementById('operationMessage').innerText = error;
				document.getElementById('operationMessage').style.color = 'red'; // Делаем текст красным
				console.error('Error:', error);
			})
			.finally(() => {
				submitButton.disabled = false; // Включаем кнопку обратно
				submitButton.innerText = 'Create';
			});
		});
		</script>
`
	ScriptOperationAccountCache string = `
		<script>
			document.getElementById('formAccountMoney').addEventListener('submit', function(event) {
			event.preventDefault(); // Предотвращаем стандартную отправку формы
		
			var submitButton = document.getElementById('accountMoneyButton');
			submitButton.disabled = true; // Отключаем кнопку
			submitButton.innerText = 'Processing...'; // Меняем текст кнопки
		
			var formData = {
				account_id: document.querySelector('input[name="account_id"]').value,
				money: document.querySelector('input[name="money"]').value
			};

			// Отправляем POST-запрос с учётом cookie
			fetch('{{.OperationRequest}}', {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json' // Указываем, что отправляем JSON
				},
				body: JSON.stringify(formData), // Преобразуем объект в JSON
				credentials: 'include' // Важно: включаем отправку и получение cookie
			})
			.then(response => {
				if (response.ok) {
					return response.json(); // Парсим JSON-ответ
				} else {
					return response.json().then(errorData => {
						throw new Error(errorData.error);
					});
				}
			})
			.then(data => {
				if (data.success) {
					window.location.href = '{{.HomePageRequest}}'; // Перенаправляем на GET-страницу
				} else {
					// Если есть ошибка, выводим сообщение
					document.getElementById('operationMessage').innerText = data.error;
				}
			})
			.catch(error => {
				document.getElementById('operationMessage').innerText = error;
				document.getElementById('operationMessage').style.color = 'red'; // Делаем текст красным
				console.error('Error:', error);
			})
			.finally(() => {
				submitButton.disabled = false; // Включаем кнопку обратно
				submitButton.innerText = 'Create';
			});
		});
		</script>
`
)
