const API_BASE_URL = 'http://localhost:8080';
const defaultHeaders = {
    'Content-Type': 'application/json',
};

let currentUser  = localStorage.getItem('currentUser') || null;
let currentUserID  = localStorage.getItem('currentUserID') || null;
let authToken = localStorage.getItem('authToken') || null;

const loginBtn = document.getElementById('login-btn');
const signupBtn = document.getElementById('signup-btn');
const logoutBtn = document.getElementById('logout-btn');
const welcomeMessage = document.getElementById('welcome-message');
const currentUserNameSpan = document.getElementById('current-user-name');
const addTransactionForm = document.getElementById('add-transaction-form');
const addTransactionBtn = document.getElementById('add-transaction-btn');
const transactionsContainer = document.getElementById('transactions-container');
const generateAmountBtn = document.getElementById('generate-amount-btn');
const modal = document.getElementById('qr-modal');
const qrCodeContainer = document.getElementById('qr-code');
const qrUrl = document.getElementById('qr-url');

async function updateUI() {
    if (authToken) {
        loginBtn.classList.add('hidden');
        signupBtn.classList.add('hidden');
        logoutBtn.classList.remove('hidden');
        addTransactionForm.classList.remove('hidden');
        welcomeMessage.classList.remove('hidden');
        currentUserNameSpan.textContent = currentUser;
        await fetchAndRenderTransactions();
    } else {
        loginBtn.classList.remove('hidden');
        signupBtn.classList.remove('hidden');
        logoutBtn.classList.add('hidden');
        welcomeMessage.classList.add('hidden');
        addTransactionForm.classList.add('hidden');
        transactionsContainer.innerHTML = '<p>Please log in to view transactions.</p>';
    }
}

async function fetchAndRenderTransactions() {
    try {
        let url = `${API_BASE_URL}/api/transactions/`;


        const response = await fetch(url, {
            method: 'GET',
            headers: {
                'Authorization': `Bearer ${authToken}`,
                'Content-Type': 'application/json',
            },
        });

        if (response.ok) {
            const transactions = await response.json();
            renderTransactions(transactions);
        } else {
            throw new Error('Failed to fetch transactions');
        }
    } catch (error) {
        console.error('Error fetching transactions:', error);
        transactionsContainer.innerHTML = '<p>Error loading transactions. Please try again later.</p>';
    }
}

async function generateRandomAmount() {
    const amountField = document.getElementById('transaction-amount');
    amountField.value = (Math.random() * (100.0 - 0.1) + 0.1).toFixed(2);
}

function renderTransactions(transactions) {
    const transactionsContainer = document.getElementById('transactions-container');
    transactionsContainer.innerHTML = transactions.map(transaction => `
        <div class="transaction">
            <p><strong>Description:</strong> ${transaction.description}</p>
            <p><strong>Amount:</strong> $${transaction.amount.toFixed(2)}</p>
            <p><strong>Status:</strong> ${transaction.order_status}</p>
            <p><strong>Created At:</strong> ${new Date(transaction.created_at).toLocaleString()}</p>
            <p><strong>Updated At:</strong> ${new Date(transaction.updated_at).toLocaleString()}</p>
            <div class="transaction-actions">
                <button class="btn delete-btn" onclick="deleteTransaction('${transaction.id}')">Delete</button>
            </div>
        </div>
    `).join('');
}

async function addTransaction() {
    const description = document.getElementById('transaction-description').value;
    const amount = parseFloat(document.getElementById('transaction-amount').value);  // Convert to a number

    // Check if amount is a valid number
    if (isNaN(amount) || !(amount && description)) {
        alert('Please enter both a valid amount and description');
        return;
    }

    try {
        const response = await fetch(`${API_BASE_URL}/api/transactions/`, {
            method: 'POST',
            headers: {
                'Authorization': `Bearer ${authToken}`,
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ amount, description }),
        });

        if (!response.ok) {
            throw new Error('Failed to add transaction');
        }

        const newTransaction = await response.json();
        console.log('New transaction added:', newTransaction);

        // Extract the redirect_url from the response
        const { redirect_url, status } = newTransaction;

        // Show a pop-up with a QR code if the status is pending
        if (status === 'pending') {
            showQRCode(redirect_url);
        }

        // Clear the form
        document.getElementById('transaction-amount').value = '';
        document.getElementById('transaction-description').value = '';

        await fetchAndRenderTransactions();
    } catch (error) {
        console.error('Error adding transaction:', error);
        alert('Failed to add transaction. Please try again later.');
    }
}

// Function to generate and show QR code
function showQRCode(url) {

    // Clear previous QR code (if any)
    qrCodeContainer.innerHTML = '';

    // Generate new QR code
    const qrCode = new QRCode(qrCodeContainer, {
        text: url,
        width: 200,
        height: 200,
    });
    qrUrl.setAttribute('href', url);

    // Show the modal
    modal.classList.remove('hidden');

    // Close button event
    document.getElementById('close-qr-btn').addEventListener('click', () => {
        modal.classList.add('hidden');
    });
}


async function deleteTransaction(id) {
    try {
        const response = await fetch(`${API_BASE_URL}/api/transactions/${id}/`, {
            method: 'DELETE',
            headers: {
                'Authorization': `Bearer ${authToken}`,
                'Content-Type': 'application/json',
            },
        });

        if (!response.ok) {
            throw new Error('Failed to delete transaction');
        }
        await fetchAndRenderTransactions();
    } catch (error) {
        console.error('Error deleting on transaction:', error);
        alert('Failed to delete transaction. Please try again later.');
    }
}

async function signup() {
    const name = prompt('Enter your name:');
    const email = prompt('Enter your email:');
    const password = prompt('Enter a password:');

    const userData = {
        name: name,
        email: email,
        password: password
    };

    fetch(`${API_BASE_URL}/auth/signup`, {
        method: 'POST',
        headers: defaultHeaders,
        body: JSON.stringify(userData)
    })
        .then(response => {
            if (response.status === 200) {
                return response.json();
            } else {
                throw new Error('Signup failed');
            }
        })
        .then(data => {
            console.log('Signup successful.');
            alert('Signup successful! Please login to continue.');
        })
        .catch(error => {
            console.error('Error:', error);
            alert('Signup failed. Please try again.');
        });
}

async function login() {
    const email = prompt('Enter your email:');
    const password = prompt('Enter your password:');

    if (email && password) {
        try {
            const response = await fetch(`${API_BASE_URL}/auth/login`, {
                method: 'POST',
                headers: defaultHeaders,
                body: JSON.stringify({ email, password }),
            });

            if (!response.ok) {
                alert('Login failed. Please try again.');
            }
            const data = await response.json();
            authToken = data.token;
            currentUser = data.name;
            currentUserID= data.id

            localStorage.setItem('authToken', authToken);
            await updateUI();

        } catch (error) {
            console.error('Error during login:', error);
            alert('An error occurred during login. Please try again.');
        }
    }
}

function logout() {
    authToken = null;
    currentUser = null;
    currentUserID = null;
    localStorage.removeItem('authToken');
    updateUI();
}

loginBtn.addEventListener('click', login);
signupBtn.addEventListener('click', signup);
logoutBtn.addEventListener('click', logout);
addTransactionBtn.addEventListener('click', addTransaction);
generateAmountBtn.addEventListener('click', generateRandomAmount);
// Initial UI update
updateUI();