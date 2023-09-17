class ArticlesAPI {
    baseURL = process.env.ARTICLES_BACKEND_URL || "http://localhost:1200/articles";

    // Create an article
    async createArticle(articleData) {
        const response = await fetch(this.baseURL, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(articleData)
        });

        if (!response.ok) {
            throw new Error(`Failed to create article: ${response.json()}`);
        }
        return await response.json();
    }

    // Get a list of articles with optional filter options
    async getArticles(filters = {}) {
        let queryParameters = new URLSearchParams();
        for (let key in filters) {
            queryParameters.append(key, filters[key]);
        }

        const response = await fetch(`${this.baseURL}?${queryParameters}`);

        if (!response.ok) {
            throw new Error(`Failed to fetch articles: ${response.json()}`);
        }
        return await response.json();
    }

    // Get a single article by its article ID
    async getArticleById(id) {
        const response = await fetch(`${this.baseURL}/${id}`);

        if (!response.ok) {
            throw new Error(`Failed to fetch article with ID ${id}: ${response.json()}`);
        }
        return await response.json();
    }

    // Update a single article by its ID
    async updateArticle(id, articleData) {
        const response = await fetch(`${this.baseURL}/${id}`, {
            method: 'UPDATE',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(articleData)
        });

        if (!response.ok) {
            throw new Error(`Failed to update article with ID ${id}: ${response.json()}`);
        }
        return await response.json();
    }

    // Delete a single article by its ID
    async deleteArticle(id) {
        const response = await fetch(`${this.baseURL}/${id}`, {
            method: 'DELETE'
        });

        if (!response.ok) {
            throw new Error(`Failed to delete article with ID ${id}: ${response.json()}`);
        }
        return await response.json();
    }
}

export default ArticlesAPI;