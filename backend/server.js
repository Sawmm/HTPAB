const express = require('express');
const fs = require('fs');
const path = require('path');
const { parse } = require('csv-parse/sync');
const { body, validationResult } = require('express-validator');
const Fuse = require('fuse.js');
const cors = require('cors');

const app = express();
const port = 5000;

// Middleware
app.use(cors());
app.use(express.json());

// DJ Performances data and Fuse instance
let djPerformances = [];
let fuse;

function loadDJList() {
    try {
        const fileContent = fs.readFileSync('berghain_lineup.csv', 'utf8');
        const records = parse(fileContent, {
            columns: true,
            skip_empty_lines: true
        });

        djPerformances = records.map(record => ({
            date: record.date,
            name: record.name.trim(),
            label: record.label.trim(),
            time: record.time,
            floor: record.floor,
            closing: record.closing === 'TRUE',
            year: parseInt(record.year)
        }));

        // Initialize Fuse with the DJ list
        const fuseOptions = {
            includeScore: true,
            threshold: 0.4,
            keys: ['name']
        };
        fuse = new Fuse(djPerformances, fuseOptions);
        
        console.log('DJ list loaded successfully');
    } catch (err) {
        console.error('Error loading or parsing CSV:', err);
    }
}

// Load DJ list initially
loadDJList();

// Reload DJ list every hour
setInterval(loadDJList, 3600000);

// Routes
app.post('/check-dj', [
    body('djName').trim().isLength({ min: 1, max: 100 }).escape()
], (req, res) => {
    const errors = validationResult(req);
    if (!errors.isEmpty()) {
        return res.status(400).json({ errors: errors.array() });
    }

    const { djName } = req.body;
    const normalizedName = djName.trim().toLowerCase();
    
    const exactMatches = djPerformances.filter(perf => perf.name.toLowerCase() === normalizedName);
    if (exactMatches.length > 0) {
        // Sort performances by date, most recent first
        const sortedPerformances = exactMatches.sort((a, b) => new Date(b.date) - new Date(a.date));
        const latestPerformance = sortedPerformances[0];
        const performanceCount = sortedPerformances.length;

        res.json({ 
            message: `${latestPerformance.name} has played at Berghain ${performanceCount} time${performanceCount > 1 ? 's' : ''}!`,
            details: `Last played on ${latestPerformance.date} at ${latestPerformance.time} on ${latestPerformance.floor}.`,
            performances: sortedPerformances
        });
    } else {
        // Perform fuzzy search
        const results = fuse.search(normalizedName);
        if (results.length > 0 && results[0].score < 0.4) {
            const suggestion = results[0].item;
            const similarMatches = djPerformances.filter(perf => perf.name.toLowerCase() === suggestion.name.toLowerCase());
            const sortedPerformances = similarMatches.sort((a, b) => new Date(b.date) - new Date(a.date));
            const performanceCount = sortedPerformances.length;

            res.json({ 
                message: `${djName} might have played at Berghain. Did you mean ${suggestion.name}?`,
                suggestion: suggestion.name,
                details: `${suggestion.name} has played ${performanceCount} time${performanceCount > 1 ? 's' : ''}, last on ${sortedPerformances[0].date} at ${sortedPerformances[0].time} on ${sortedPerformances[0].floor}.`,
                performances: sortedPerformances
            });
        } else {
            res.json({ message: `${djName} has not played at Berghain (according to our records).` });
        }
    }
});

// Error handling middleware
app.use((err, req, res, next) => {
    console.error(err.stack);
    res.status(500).json({ message: 'Something went wrong!' });
});

// Start the server
app.listen(port, () => {
    console.log(`Server running at http://localhost:${port}`);
});