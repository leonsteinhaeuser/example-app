import React, { useState } from 'react';
import '../styles/create_article_form.css';

// Assuming you have an external WYSIWYG library installed, e.g. 'react-quill'
import ReactQuill from 'react-quill';
import 'react-quill/dist/quill.snow.css'; // Import styles

const CreateArticle = ({ authors, articleAPI }) => {
    const [formData, setFormData] = useState({
        title: '',
        description: '',
        content: '',
        published: false,
        tags: [''],
        author: '',
        co_authors: [''],
    });

    const handleInputChange = (event) => {
        const { name, value, type, checked } = event.target;
        if (type === "checkbox") {
            setFormData((prevData) => ({ ...prevData, [name]: checked }));
        } else {
            setFormData((prevData) => ({ ...prevData, [name]: value }));
        }
    };

    const handleContentChange = (value) => {
        setFormData((prevData) => ({ ...prevData, content: value }));
    };

    const handleTagChange = (index, value) => {
        const newTags = [...formData.tags];
        newTags[index] = value;
        setFormData((prevData) => ({ ...prevData, tags: newTags }));
    };

    const addTag = () => {
        setFormData((prevData) => ({ ...prevData, tags: [...prevData.tags, ''] }));
    };

    const removeTag = (index) => {
        const newTags = [...formData.tags];
        newTags.splice(index, 1);
        setFormData((prevData) => ({ ...prevData, tags: newTags }));
    };

    const handleSubmit = async (event) => {
        event.preventDefault();
        try {
            articleAPI.createArticle(formData);
        } catch (error) {
            console.error('Error submitting data:', error);
        }
    };

    return (
        <form onSubmit={handleSubmit}>
            <input
                type="text"
                name="title"
                placeholder="Title"
                value={formData.title}
                onChange={handleInputChange}
            />
            <textarea
                name="description"
                placeholder="Description"
                value={formData.description}
                onChange={handleInputChange}
            />
            <ReactQuill
                value={formData.content}
                onChange={handleContentChange}
            />
            <label>
                Published:
                <input
                    type="checkbox"
                    name="published"
                    checked={formData.published}
                    onChange={handleInputChange}
                />
            </label>

            <div>
                {formData.tags.map((tag, index) => (
                    <div key={index}>
                        <input
                            type="text"
                            placeholder={`Tag ${index + 1}`}
                            value={tag}
                            onChange={(e) => handleTagChange(index, e.target.value)}
                        />
                        <button onClick={() => removeTag(index)}>Remove</button>
                    </div>
                ))}
                <button onClick={addTag}>Add Tag</button>
            </div>

            <select name="author" value={formData.author} onChange={handleInputChange}>
                <option value="" disabled>Select Author</option>
                {authors.map((author, index) => (
                    <option key={index} value={author}>{author}</option>
                ))}
            </select>

            <select multiple name="co_authors" value={formData.co_authors} onChange={(e) => {
                const values = [...e.target.selectedOptions].map(option => option.value);
                setFormData(prevData => ({ ...prevData, co_authors: values }));
            }}>
                {authors.map((author, index) => (
                    <option key={index} value={author}>{author}</option>
                ))}
            </select>

            <button type="submit">Submit</button>
        </form>
    );
};

export default CreateArticle;
