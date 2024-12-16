import React, { useEffect, useState } from 'react';
import { GetSettings, UpdateSettings } from '../../wailsjs/go/main/App';
import { contract } from '../../wailsjs/go/models';

interface SettingsFormProps {
    onSubmit: (settings: contract.Settings) => void;
    initialSettings?: contract.Settings;
}

const SettingsForm: React.FC<SettingsFormProps> = ({ onSubmit }) => {
    const [settings, setSettings] = useState<contract.Settings>({ Username: "", ExcludeDraft: false, ExcludeConflicting: false });
    useEffect(() => {
        const fetchSettings = async () => {
            const settings = await GetSettings()
            setSettings(settings)
        }
        fetchSettings()
    }, [])
    const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        if (e.target.name == 'Username') {
            setSettings({ ...settings, [e.target.name]: e.target.value });
        } else {
            setSettings({ ...settings, [e.target.name]: e.target.checked });
        }
    }

    const handleSubmit = (e: React.FormEvent) => {
        e.preventDefault();
        UpdateSettings(settings).catch(alert)
    };

    return (
        <form onSubmit={handleSubmit} className="container">
            <input name="Username" data-testid="username" placeholder="GitHub Username" value={settings.Username} onChange={handleChange} required />
            <input type="checkbox" name="ExcludeDraft" data-testid="exclude-draft" checked={settings.ExcludeDraft} onChange={handleChange} required />
            <input type="checkbox" name="ExcludeConflicting" data-testid="exclude-conflicting" checked={settings.ExcludeConflicting} onChange={handleChange} required />
            <button type="submit">Save</button>
        </form>
    );
};

export default SettingsForm;
