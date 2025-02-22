// Copyright 2023 The Perses Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

import { Box, Grid, Stack } from '@mui/material';
import { useMemo, useState } from 'react';
import HomeIcon from 'mdi-material-ui/Home';
import { useNavigate } from 'react-router-dom';
import { DashboardSelector, ProjectResource } from '@perses-dev/core';
import { useProjectList } from '../../model/project-client';
import { CreateProjectDialog, CreateDashboardDialog } from '../../components/dialogs';
import { StackCrumb, TitleCrumb } from '../../components/breadcrumbs/breadcrumbs';
import { useIsMobileSize } from '../../utils/browser-size';
import { CRUDButton } from '../../components/CRUDButton/CRUDButton';
import { InformationSection } from './InformationSection';
import { RecentDashboards } from './RecentDashboards';
import { ProjectsAndDashboards } from './ProjectsAndDashboards';
import { ImportantDashboards } from './ImportantDashboards';

function HomeView() {
  // Navigate to the project page if the project has been successfully added
  const navigate = useNavigate();
  const isMobileSize = useIsMobileSize();

  const { data } = useProjectList();

  const projectOptions = useMemo(() => {
    return (data || []).map((project) => project.metadata.name); // TODO: remove projects without create dashboard perm
  }, [data]);

  const handleAddProjectDialogSubmit = (entity: ProjectResource) => navigate(`/projects/${entity.metadata.name}`);
  const handleAddDashboardDialogSubmit = (dashboardSelector: DashboardSelector) =>
    navigate(`/projects/${dashboardSelector.project}/dashboard/new`, { state: dashboardSelector.dashboard });

  // Open/Close management for dialogs
  const [isAddProjectDialogOpen, setIsAddProjectDialogOpen] = useState(false);
  const [isAddDashboardDialogOpen, setIsAddDashboardDialogOpen] = useState(false);

  const handleAddProjectDialogOpen = () => {
    setIsAddProjectDialogOpen(true);
  };
  const handleAddProjectDialogClose = () => {
    setIsAddProjectDialogOpen(false);
  };
  const handleAddDashboardDialogOpen = () => {
    setIsAddDashboardDialogOpen(true);
  };
  const handleAddDashboardDialogClose = () => {
    setIsAddDashboardDialogOpen(false);
  };

  return (
    <Stack sx={{ width: '100%', overflowX: 'hidden' }} m={isMobileSize ? 1 : 2} gap={1}>
      <Box sx={{ width: '100%' }}>
        <Stack direction="row" alignItems="center" justifyContent="space-between">
          <StackCrumb>
            <HomeIcon fontSize={'large'} />
            <TitleCrumb>Home</TitleCrumb>
          </StackCrumb>
          <Stack direction="row" gap={isMobileSize ? 0.5 : 2}>
            <CRUDButton action="create" scope="Project" variant="contained" onClick={handleAddProjectDialogOpen}>
              Add Project
            </CRUDButton>
            <CRUDButton
              variant="contained"
              onClick={handleAddDashboardDialogOpen}
              disabled={projectOptions.length === 0}
            >
              Add Dashboard
            </CRUDButton>
            <CreateProjectDialog
              open={isAddProjectDialogOpen}
              onClose={handleAddProjectDialogClose}
              onSuccess={handleAddProjectDialogSubmit}
            />
            <CreateDashboardDialog
              open={isAddDashboardDialogOpen}
              projectOptions={projectOptions}
              onClose={handleAddDashboardDialogClose}
              onSuccess={handleAddDashboardDialogSubmit}
            />
          </Stack>
        </Stack>
      </Box>
      <Grid container columnSpacing={8}>
        <Grid item xs={12} lg={8}>
          <RecentDashboards />
          <ProjectsAndDashboards />
        </Grid>
        <Grid item xs={12} lg={4}>
          <ImportantDashboards />
          <InformationSection />
        </Grid>
      </Grid>
    </Stack>
  );
}

export default HomeView;
