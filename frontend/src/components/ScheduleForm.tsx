import React, { useState } from 'react';

interface ScheduleFormProps {
  onSubmit: (schedule: string) => void;
  initialSchedule?: string;
}

const ScheduleForm: React.FC<ScheduleFormProps> = ({ onSubmit, initialSchedule }) => {
  const [schedule, setSchedule] = useState<string>(initialSchedule || '');

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    onSubmit(schedule);
  };

  return (
    <form onSubmit={handleSubmit} className='container'>
      <input
        type="text"
        placeholder="Cron-like Schedule"
        value={schedule}
        onChange={(e) => setSchedule(e.target.value)}
        required
      />
      <button type="submit">Set Schedule</button>
    </form>
  );
};

export default ScheduleForm;
