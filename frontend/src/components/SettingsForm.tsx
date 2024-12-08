import React, { useState } from 'react';

interface Settings {
    username: string;
}

interface SettingsFormProps {
    onSubmit: (settings: Settings) => void;
    initialSettings?: Settings;
}

const SettingsForm: React.FC<SettingsFormProps> = ({ onSubmit, initialSettings }) => {
    const [settings, setSettings] = useState<Settings>(initialSettings || { username: '' });

    const handleChange = (e: React.ChangeEvent<HTMLInputElement>) =>
        setSettings({ ...settings, [e.target.name]: e.target.value });

    const handleSubmit = (e: React.FormEvent) => {
        e.preventDefault();
        onSubmit(settings);
    };

    return (
        <form onSubmit={handleSubmit} className="container">
            <input name="username" placeholder="GitHub Username" value={settings.username} onChange={handleChange} required />
            <button type="submit">Save</button>
        </form>
    );
};

export default SettingsForm;
