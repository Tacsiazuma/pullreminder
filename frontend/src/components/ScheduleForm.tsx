import React, { useState } from 'react';

interface ScheduleFormProps {
  onSubmit: (cronExpression: string) => void;
}

const ScheduleForm: React.FC<ScheduleFormProps> = ({ onSubmit }) => {
  const [scheduleOption, setScheduleOption] = useState<string>('morning');

  const scheduleOptions = [
    { label: 'Every Morning (9 AM)', value: 'morning' },
    { label: 'Every Hour (During Working Hours: 9 AM - 5 PM)', value: 'working_hours' },
    { label: 'Every Day at Midnight', value: 'midnight' },
    { label: 'Every Hour', value: 'hourly' },
    { label: 'Every Monday Morning (9 AM)', value: 'monday_morning' },
  ];

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    const cronExpression = convertToCron(scheduleOption);
    onSubmit(cronExpression);
  };

  const convertToCron = (option: string): string => {
    switch (option) {
      case 'morning':
        return '0 9 * * *'; // At 9:00 AM every day
      case 'working_hours':
        return '0 9-17 * * *'; // At the top of every hour from 9 AM to 5 PM
      case 'midnight':
        return '0 0 * * *'; // At midnight every day
      case 'hourly':
        return '0 * * * *'; // At the top of every hour
      case 'monday_morning':
        return '0 9 * * 1'; // At 9:00 AM every Monday
      default:
        return '';
    }
  };

  return (
    <form onSubmit={handleSubmit} className="schedule-container">
      <h2>Select Schedule</h2>
      <div>
        {scheduleOptions.map((option) => (
          <label key={option.value} style={{ display: 'block', margin: '10px 0' }}>
            <input
              type="radio"
              name="schedule"
              value={option.value}
              checked={scheduleOption === option.value}
              onChange={() => setScheduleOption(option.value)}
            />
            {option.label}
          </label>
        ))}
      </div>
      <button type="submit">Set Schedule</button>
    </form>
  );
};

export default ScheduleForm;
