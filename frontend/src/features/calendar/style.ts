import styled from '@emotion/styled'

export const StyleWrapper = styled.div`
    .fc .fc-toolbar.fc-header-toolbar {
        margin-bottom: 0;
    }

    .fc .fc-col-header-cell {
        font-size: 0.75rem;
        font-weight: normal;
        color: #b6b5b3;
        border: none;
    }

    .fc .fc-toolbar-title {
        font-size: 1rem;
        color: #37362f;
    }

    .fc .fc-button-primary {
        font-size: 0.75rem;
        border: none;
        outline: none;
    }

    .fc .fc-button-primary:not(:disabled):focus,
    .fc .fc-button-primary:not(:disabled).fc-button-focus {
        box-shadow: none;
    }

    .fc .fc-button-primary:active {
        border: none;
        outline: none;
    }

    .fc .fc-prev-button {
        background-color: #ffffff00;
        color: #acaba9;
    }

    .fc .fc-next-button {
        background-color: #ffffff00;
        color: #acaba9;
    }

    .fc .fc-prev-button:hover,
    .fc .fc-next-button:hover {
        background-color: #f0f0f0;
        color: #acaba9;
    }

    .fc .fc-prev-button:active,
    .fc .fc-next-button:active {
        background-color: #f0f0f0;
        color: #acaba9;
    }

    .fc .fc-scrollgrid {
        border-width: 0;
    }

    .fc .fc-scrollgrid-section > * {
        border: none;
    }

    .fc .fc-scrollgrid-sync-table {
        border: 1px;
    }
    
    .fc .fc-daygrid-day-number {
        font-size: 0.75rem;
    }

    .fc .fc-day-today {
        background-color: #ffffff00;
    }

    .fc .fc-timegrid-axis-cushion {
        color: #acaba9;
    }
    
    .fc .fc-timegrid-slot-label {
        font-size: 0.75rem;
        color: #acaba9;
    }
`;

export const CircleNumber = styled.div`
  display: inline-flex;
  justify-content: center;
  align-items: center;
  border-radius: 50%;
  flex-flow: column;
  vertical-align: top;
  background: #eb5757;
  color: white;
  width: 1rem; /* 16px */
  height: 1rem; /* 16px */

  @media (min-width: 768px) {
    width: 1.5rem; /* 24px */
    height: 1.5rem; /* 24px */
  }

  @media (min-width: 1024px) {
    width: 2rem; /* 32px */
    height: 2rem; /* 32px */
`

export const CircleToday = styled.div`
  display: inline-flex;
  justify-content: center;
  align-items: center;
  border-radius: 10%;
  flex-flow: column;
  vertical-align: top;
  background: #eb5757;
  color: white;
  width: 2.5rem; /* 60px */
  height: 0.75rem; /* 20px */

  @media (min-width: 768px) {
    width: 3rem; /* 80px */
    height: 1rem; /* 16px */
  }

  @media (min-width: 1024px) {
    width: 3.75rem; /* 60px */
    height: 1.25rem; /* 20px */
  }
`