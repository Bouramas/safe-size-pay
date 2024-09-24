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
const addMovieForm = document.getElementById('add-movie-form');
const addTransactionBtn = document.getElementById('add-transaction-btn');
const moviesContainer = document.getElementById('movies-container');
const sortSelect = document.getElementById('sort-select');
const sortOptions = document.getElementById('sort-options');
const generateAmountBtn = document.getElementById('generate-amount-btn');

async function updateUI() {
    if (authToken) {
        loginBtn.classList.add('hidden');
        signupBtn.classList.add('hidden');
        sortOptions.classList.remove('hidden');
        logoutBtn.classList.remove('hidden');
        addMovieForm.classList.remove('hidden');
        welcomeMessage.classList.remove('hidden');
        currentUserNameSpan.textContent = currentUser;
        await fetchAndRenderMovies();
    } else {
        loginBtn.classList.remove('hidden');
        signupBtn.classList.remove('hidden');
        logoutBtn.classList.add('hidden');
        sortOptions.classList.add('hidden');
        welcomeMessage.classList.add('hidden');
        addMovieForm.classList.add('hidden');
        moviesContainer.innerHTML = '<p>Please log in to view movies.</p>';
    }
}

async function fetchAndRenderMovies(userId = null) {
    try {
        let url = `${API_BASE_URL}/api/movies/`;
        if (userId) {
            url = `${API_BASE_URL}/api/users/${userId}/movies`;
        }

        const response = await fetch(url, {
            method: 'GET',
            headers: {
                'Authorization': `Bearer ${authToken}`,
                'Content-Type': 'application/json',
            },
        });

        if (response.ok) {
            const movies = await response.json();
            renderMovies(movies);
        } else {
            throw new Error('Failed to fetch movies');
        }
    } catch (error) {
        console.error('Error fetching movies:', error);
        moviesContainer.innerHTML = '<p>Error loading movies. Please try again later.</p>';
    }
}

async function generateRandomAmount() {
    const amountField = document.getElementById('transaction-amount');
    amountField.value = (Math.random() * (100.0 - 0.1) + 0.1).toFixed(2);
}

function renderMovies(movies) {
    const sortedMovies = sortMovies(movies);
    moviesContainer.innerHTML = sortedMovies.map(movie => `
        <div class="movie">
            <h3>${movie.title}</h3>
            <p>${movie.description}</p>
            <p>Added by: <a href="#" onclick="fetchAndRenderMovies('${movie.user_id}')">${movie.user_name}</a> on ${new Date(movie.date_added).toLocaleString()}</p>
            <div class="movie-actions">
                <div class="vote-counts">
                    <span>Likes: <span class="like-count">${movie.like_votes}</span></span>
                    <span>Hates: <span class="hate-count">${movie.hate_votes}</span></span>
                </div>
                ${movie.user_id !== currentUserID ? `
                    <div class="vote-buttons">
                        ${movie.vote_id > 0 ? `
                            ${movie.vote_type === 'hate' ? `
                                        <span>You hate this movie</span>
                                        <button onclick="updateVote(${movie.vote_id}, 'like')">Like</button>
                                    ` : `
                                        <span>You like this movie</span>
                                        <button onclick="updateVote(${movie.vote_id}, 'hate')">Hate</button>
                                    `}
                            <button onclick="deleteVote(${movie.vote_id})">Delete</button>
                        ` : `
                            <button onclick="voteMovie(${movie.id}, 'like')">Like</button>
                            <button onclick="voteMovie(${movie.id}, 'hate')">Hate</button>
                        `}
                    </div>
                ` : ''}
            </div>
        </div>
    `).join('');
}

function sortMovies(movies) {
    return movies.sort((a, b) => {
        switch (sortSelect.value) {
            case 'likes':
                return b.like_votes - a.like_votes;
            case 'hates':
                return b.hate_votes - a.hate_votes;
            default:
                return new Date(b.date_added) - new Date(a.date_added);
        }
    });
}

async function addTransaction() {
    const title = document.getElementById('movie-title').value;
    const description = document.getElementById('transaction-description').value;
    const amount = document.getElementById('transaction-amount').value;
    if (!(title && description)) {
        alert('Please enter both title and description');
    }
    try {
        const response = await fetch(`${API_BASE_URL}/api/movies/`, {
            method: 'POST',
            headers: {
                'Authorization': `Bearer ${authToken}`,
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({title, description}),
        });

        if (!response.ok) {
            throw new Error('Failed to add movie');
        }
        const newMovie = await response.json();
        console.log('New movie added:', newMovie);

        // Clear the form
        document.getElementById('movie-title').value = '';
        document.getElementById('movie-description').value = '';

        // Fetch and render movies again
        await fetchAndRenderMovies();

        alert('Movie added successfully!');

    } catch (error) {
        console.error('Error adding movie:', error);
        alert('Failed to add movie. Please try again later.');
    }

}

// noinspection JSUnusedGlobalSymbols
async function voteMovie(id, voteType) {
    try {
        const response = await fetch(`${API_BASE_URL}/api/movies/${id}/votes`, {
            method: 'POST',
            headers: {
                'Authorization': `Bearer ${authToken}`,
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ vote_type: voteType }),
        });

        if (!response.ok) {
            throw new Error('Failed to vote on movie');
        }
        await fetchAndRenderMovies();
    } catch (error) {
        console.error('Error voting on movie:', error);
        alert('Failed to vote on movie. Please try again later.');
    }
}

// noinspection JSUnusedGlobalSymbols
async function updateVote(id, voteType) {
    try {
        const response = await fetch(`${API_BASE_URL}/api/votes/${id}/`, {
            method: 'PATCH',
            headers: {
                'Authorization': `Bearer ${authToken}`,
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ vote_type: voteType }),
        });

        if (!response.ok) {
            throw new Error('Failed to update vote');
        }
        await fetchAndRenderMovies();
    } catch (error) {
        console.error('Error voting on movie:', error);
        alert('Failed to vote on movie. Please try again later.');
    }
}

async function deleteVote(id) {
    try {
        const response = await fetch(`${API_BASE_URL}/api/votes/${id}/`, {
            method: 'DELETE',
            headers: {
                'Authorization': `Bearer ${authToken}`,
                'Content-Type': 'application/json',
            },
        });

        if (!response.ok) {
            throw new Error('Failed to delete vote');
        }
        await fetchAndRenderMovies();
    } catch (error) {
        console.error('Error voting on movie:', error);
        alert('Failed to vote on movie. Please try again later.');
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
sortSelect.addEventListener('change', () => fetchAndRenderMovies());
generateAmountBtn.addEventListener('click', generateRandomAmount);
// Initial UI update
updateUI();